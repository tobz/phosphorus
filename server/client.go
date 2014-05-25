package server

import (
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"time"

	"github.com/rcrowley/go-metrics"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/database/models"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
	"github.com/tobz/phosphorus/statistics"
)

var globalBytesSent metrics.Counter
var globalBytesReceived metrics.Counter

func init() {
	globalBytesSent = metrics.GetOrRegisterCounter("client.bytesSent", statistics.Registry)
	globalBytesReceived = metrics.GetOrRegisterCounter("client.bytesReceived", statistics.Registry)
}

type Client struct {
	errors chan error
	logger *log.Logger

	server     *Server
	connection *net.TCPConn

	inbound  *network.PacketReader
	outbound *network.PacketWriter

	sessionId    uint16
	connectionId uint32

	version constants.ClientVersion
	state   constants.ClientState

	lastPingTime time.Time

	account *models.Account
}

func NewClient(server *Server, connection *net.TCPConn, connectionId uint32) *Client {
	c := &Client{
		server:     server,
		connection: connection,

		inbound:  network.NewPacketReader(connection, globalBytesReceived),
		outbound: network.NewPacketWriter(connection, globalBytesSent),

		connectionId: connectionId,
	}

	clientLogger := log.NewLogger()
	clientLogger.SetPrefixer(func(s string) string {
		var clientPrefix string
		if c.Account() != nil {
			clientPrefix = fmt.Sprintf("[%s / %s]", c.connection.RemoteAddr().String(), c.Account().Username)
		} else {
			clientPrefix = fmt.Sprintf("[%s]", c.connection.RemoteAddr().String())
		}

		return fmt.Sprintf("%s %s", clientPrefix, s)
	})

	c.logger = clientLogger

	return c
}

func (c *Client) Start() {
	// Set our client state.
	c.state = constants.ClientStateConnecting

	// Start handling our network connection.
	go func() {
		// Zero tolerance error handling
		defer func() {
			err := recover()
			if err != nil {
				if err == io.EOF {
					c.logger.Info("client", "Connection closed from client side.")
				} else {
					c.logger.Error("client", "Caught an error while in network loop: %s", err)
					c.logger.Error("client", "Stack trace: %s", debug.Stack())
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

func (c *Client) Logger() interfaces.Logger {
	return c.logger
}

func (c *Client) SetAccount(account *models.Account) {
	c.account = account
}

func (c *Client) Account() *models.Account {
	return c.account
}

func (c *Client) SetClientVersion(version constants.ClientVersion) {
	c.version = version
}

func (c *Client) ClientVersion() constants.ClientVersion {
	return c.version
}

func (c *Client) SetClientState(state constants.ClientState) {
	c.state = state
}

func (c *Client) ClientState() constants.ClientState {
	return c.state
}

func (c *Client) LastPingTime() time.Time {
	return c.lastPingTime
}

func (c *Client) MarkPingTime() {
	c.lastPingTime = time.Now()
}

func (c *Client) ConnectionID() uint32 {
	return c.connectionId
}

func (c *Client) SetSessionID(id uint16) {
	c.sessionId = id
}

func (c *Client) SessionID() uint16 {
	return c.sessionId
}

func (c *Client) Server() interfaces.Server {
	return c.server
}

func (c *Client) Send(p interfaces.Packet) error {
	packetType := "TCP"
	if p.Type() == constants.PacketUDP {
		packetType = "UDP"
	}

	c.logger.Debug("client", "Sending packet %s(0x%X) -> %d bytes", packetType, uint8(p.Code()), p.Len())

	return c.outbound.Write(p)
}
