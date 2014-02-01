package phosphorus

import "fmt"
import "net"
import "hash/crc32"
import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/utils"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/managers"
import "github.com/tobz/phosphorus/network"

type Client struct {
	coordinator *utils.Coordinator
	errors      chan error

	server     *Server
	connection *net.TCPConn

	bufPosition int
	readBuffer  []byte

	sendQueue    chan interfaces.Packet
	receiveQueue chan *network.InboundPacket

	clientId uint16
	account  interfaces.Account
}

func NewClient(server *Server, connection *net.TCPConn) *Client {
	return &Client{
		coordinator:  utils.NewCoordinator(),
		errors:       make(chan error, 1),
		server:       server,
		connection:   connection,
		readBuffer:   make([]byte, 32768),
		sendQueue:    make(chan interfaces.Packet, 128),
		receiveQueue: make(chan *network.InboundPacket, 128),
		clientId:     0,
		account:      nil,
	}
}

func (c *Client) Start() {
	// Handle any errors that crop up.
	go func() {
		stop := c.coordinator.Register()

		for {
			select {
			case <-stop:
				break
			case err := <-c.errors:
				// Log this error and stop the client.
				log.Server.Error("client", log.WrapForClient(c, "caught an error: %s", err))

				c.Stop()
				break
			}
		}
	}()

	// Start handling our network connection.
	go func() {
		stop := c.coordinator.Register()

		for {
			select {
			case <-stop:
				break
			default:
				// Make sure we have runway to receive.
				if c.bufPosition >= len(c.readBuffer) {
					c.errors <- fmt.Errorf("overflowed receive buffer: offset %d with buf size %d", c.bufPosition, len(c.readBuffer))
					break
				}

				// Read from our connection.
				n, err := c.connection.Read(c.readBuffer[c.bufPosition:])
				if err != nil {
					c.errors <- err
					break
				}

				c.bufPosition += n

				// See if we have a full packet yet.
				packet, err := c.tryForPacket()
				if err != nil {
					c.errors <- err
					break
				}

				// Stick it in the queue.
				c.receiveQueue <- packet
			}
		}
	}()

	// Start listening to our packet queues.
	go func() {
		stop := c.coordinator.Register()

		for {
			select {
			case <-stop:
				break
			case packet := <-c.receiveQueue:
				err := c.handlePacket(packet)
				if err != nil {
					c.errors <- err
					break
				}
			case packet := <-c.sendQueue:
				err := c.sendPacket(packet)
				if err != nil {
					c.errors <- err
					break
				}
			}
		}
	}()
}

func (c *Client) tryForPacket() (*network.InboundPacket, error) {
	return nil, nil
}

func (c *Client) handlePacket(packet *network.InboundPacket) error {
	return managers.DefaultPacketManager.HandlePacket(c, packet)
}

func (c *Client) Send(packet interfaces.Packet) error {
	c.sendQueue <- packet

	return nil
}

func (c *Client) sendPacket(packet interfaces.Packet) error {
	// Figure out if we have to hand this over to the server to send over UDP.
	if packet.Type() == constants.PacketType_UDP {
		return c.server.SendUDP(c, packet)
	}

	// No UDP, so just send this over our TCP connection.
	n, err := c.connection.Write(packet.Buffer())
	if err != nil {
		return err
	}

	// Make sure we sent it all.
	if n != len(packet.Buffer()) {
		return fmt.Errorf("tried to send packet with %d bytes, but only sent %d bytes", len(packet.Buffer()), n)
	}

	return nil
}

func (c *Client) GetUniqueIdentifier() uint32 {
	data := []byte(c.connection.RemoteAddr().String())
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

func (c *Client) SetAccount(account interfaces.Account) {
	c.account = account
}

func (c *Client) Account() interfaces.Account {
	return c.account
}

func (c *Client) RemoteAddr() net.Addr {
	return c.connection.RemoteAddr()
}

func (c *Client) Cleanup() {
}

func (c *Client) Stop() {
	c.coordinator.Ping()
	c.coordinator.Stop()

	// Close our socket.
	c.connection.Close()

	// Clean ourselves up.
	c.Cleanup()

	// Inform the server we're closing up shop so it can clean us up.
	c.server.RemoveClient(c)
}
