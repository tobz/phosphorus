package phosphorus

import "net"
import "github.com/tobz/phosphorus/managers"
import "github.com/tobz/phosphorus/network"

type Client struct {
	coordinator *Coordinator
	errors      chan error

	server     *Server
	connection *net.TCPConn

	bufPosition int
	readBuffer  []byte

	sendQueue    chan *network.OutboundPacket
	receiveQueue chan *network.InboundPacket

	clientId uint16
}

func NewClient(server *Server, connection *net.TCPConn) *Client {
	return &Client{
		coordinator:  NewCoordinator(),
		errors:       make(chan error, 1),
		server:       server,
		connection:   connection,
		readBuffer:   make([]byte, 32768),
		sendQueue:    make(chan *network.OutboundPacket, 128),
		receiveQueue: make(chan *network.InboundPacket, 128),
		clientId:     0,
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
					c.errors <- ClientErrorf(c, "overflowed receive buffer: offset %d with buf size %d", c.bufPosition, len(c.readBuffer))
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

func (c *Client) Send(packet *network.OutboundPacket) {
	c.sendQueue <- packet
}

func (c *Client) sendPacket(packet *network.OutboundPacket) error {
	// Figure out if we have to hand this over to the server to send over UDP.
	if packet.Type == network.PacketType_UDP {
		return c.server.SendUDP(client, packet)
	}

	// No UDP, so just send this over our TCP connection.
	n, err := c.connection.Write(packet.Buffer())
	if err != nil {
		return err
	}

	// Make sure we sent it all.
	if n != len(packet.Buffer()) {
		return ClientErrorf(c, "tried to send packet with %d bytes, but only sent %d bytes", len(packet.Buffer()), n)
	}

	return nil
}

func (c *Client) Stop() {
	c.coordinator.Ping()
	c.coordinator.Stop()
}
