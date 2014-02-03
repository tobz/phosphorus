package phosphorus

import (
    "net"

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

    ClientId uint16
    ConnectionId uint32
    Account *models.Account
}

func NewClient(server *Server, connection *net.TCPConn) *Client {
    return &Client{
        server:     server,
        connection: connection,

        inbound:  network.NewPacketReader(connection),
        outbound: network.NewPacketWriter(connection),

        ClientId:      0,
        ConnectionId: 0,
        Account: nil,
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
                log.Server.Error("client", log.WrapForClient(c, "caught an error: %s", err))
            }
        }()

        for {
            p, err := c.inbound.Next()
            if err != nil {
                panic(err)
            }

            err = c.handle(p)
            if err != nil {
                panic(err)
            }
        }
    }()
}

func (c *Client) handle(packet *network.InboundPacket) error {
    return handlers.Handle(c, packet)
}

func (c *Client) Send(packet network.Packet) error {
    return c.outbound.Write(packet)
}

func (c *Client) RemoteAddr() net.Addr {
    return c.connection.RemoteAddr()
}

func (c *Client) cleanup() {
}

func (c *Client) Stop() {
    // Inform the server we're closing up shop so it can clean us up.
    c.server.RemoveClient(c)

    // close reader/writer
    c.reader.Stop()
    c.writer.Stop()

    // close our socket.
    c.connection.Close()

    // clean ourselves up.
    c.cleanup()
}
