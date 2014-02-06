package server

import (
	"io"
	"net"

    "github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

type Client struct {
	errors chan error

	server     *Server
	connection *net.TCPConn

	inbound  *network.PacketReader
	outbound *network.PacketWriter

	clientId     uint16
    version constants.ClientVersion
	connectionId uint32
	account      interfaces.Account
}

func NewClient(server *Server, connection *net.TCPConn, connectionId uint32) *Client {
	return &Client{
		server:     server,
		connection: connection,

		inbound:  network.NewPacketReader(connection),
		outbound: network.NewPacketWriter(connection),

		clientId:     0,
		connectionId: connectionId,
		account:      nil,
	}
}

func (c *Client) Start() {
	// Start handling our network connection.
	go func() {
		// Zero tolerance error handling
		defer func() {
			err := recover()
			if err != nil {
                if err == io.EOF {
                    log.Server.ClientInfo(c, "client", "Connection closed from client side.")
                } else {
				    log.Server.ClientError(c, "client", "Caught an error while in network loop: %s", err)
                }

				c.Stop()
			}
		}()

		for {
			p, err := c.inbound.Next()
			if err != nil {
				panic(err)
			}

			if p != nil {
				err = handlers.Handle(c, p)
				if err != nil {
					panic(err)
				}
			}
		}
	}()
}

func (c *Client) cleanup() {
}

func (c *Client) Stop() {
	// Inform the server we're closing up shop so it can clean us up.
	c.server.RemoveClient(c)

	// close our socket.
	c.connection.Close()

	// clean ourselves up.
	c.cleanup()
}

// Methods to satisfy interfaces.Client
func (c *Client) Connection() interfaces.LimitedConnection {
    return c.connection
}

func (c *Client) SetAccount(account interfaces.Account) {
	c.account = account
}

func (c *Client) Account() interfaces.Account {
	return c.account
}

func (c *Client) SetClientVersion(version constants.ClientVersion) {
	c.version = version
}

func (c *Client) ClientVersion() constants.ClientVersion {
	return c.version
}

func (c *Client) ConnectionId() uint32 {
	return c.connectionId
}

func (c *Client) Server() interfaces.Server {
	return c.server
}

func (c *Client) Send(p interfaces.Packet) error {
    packetType := "TCP"
    if p.Type() == constants.PacketUDP {
        packetType = "UDP"
    }

    log.Server.ClientDebug(c, "client", "Sending packet %s(0x%X) -> %d bytes", packetType, uint8(p.Code()), p.Len())

	return c.outbound.Write(p)
}
