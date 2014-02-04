package phosphorus

import (
	"fmt"
	"io"
	"net"

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
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Server.Error("client", wrapForClient(c, "caught an error: %s", err))
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

func wrapForClient(c *Client, format string, args ...interface{}) string {
	var prefix string

	if c.Account() != nil {
		prefix = fmt.Sprintf("[%s / %s] ", c.connection.RemoteAddr().String(), c.Account().Name())
	} else {
		prefix = fmt.Sprintf("[%s] ", c.connection.RemoteAddr().String())
	}

	return fmt.Sprintf(prefix+format, args...)
}

// Methods to satisfy interfaces.Client
func (c *Client) SetAccount(account interfaces.Account) {
	c.account = account
}

func (c *Client) Account() interfaces.Account {
	return c.account
}

func (c *Client) ConnectionId() uint32 {
	return c.connectionId
}

func (c *Client) Server() interfaces.Server {
	return c.server
}

func (c *Client) Send(packet interfaces.Packet) error {
	return c.outbound.Write(packet)
}
