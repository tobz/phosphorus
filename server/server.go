package phosphorus

import (
    "net"
    "sync"

    "github.com/tobz/phosphorus/log"
    "github.com/tobz/phosphorus/network"
    "github.com/tobz/phosphorus/utils"
)

type Server struct {
    Config *ServerConfig

    tcpListener *net.TCPListener
    udpListener *net.UDPConn

    register   chan *Client
    unregister chan *Client
    stop       chan struct{}
}

func NewServer(config *ServerConfig) *Server {
    return &Server{
        Config: config,

        tcpListener: nil,
        udpListener: nil,

        clients: make(map[uint32]*Client, 128),

        register:   make(chan *Client),
        unregister: make(chan *Client),
        stop:       make(chan struct{}),
    }
}

func (s *Server) Start() error {
    log.Server.Info("server", "Starting the server...")

    // Do a bunch of random shit here.

    // Now start listening.
    tl, err := net.ListenTCP(s.config.tcpListenAddr.Network(), s.config.tcpListenAddr)
    if err != nil {
        s.Stop()
        return err
    }

    s.tcpListener = tl

    ul, err := net.ListenUDP(s.config.udpListenAddr.Network(), s.config.udpListenAddr)
    if err != nil {
        s.Stop()
        return err
    }

    s.udpListener = ul

    // Now run our workers that pull things off the wire and accept connections.
    go func() {
        for {
            select {
            case <-s.stop:
                for _, c := range s.clients {
                    c.Stop()
                }
                return
            case c := <-s.register:
                log.Server.Info("server", "Accepting new connection from client %s", c.conn.RemoteAddr().String())
                s.clients[c.ConnectionId] = c
            case c := <-s.unregister:
                log.Server.Info("server", "Closing connection for client %s", c.conn.RemoteAddr().String())
                delete(s.clients, c.ConnectionId)
            }
        }
    }()

    go func() {
        log.Server.Info("server", "Now listening for TCP connections at %s", s.config.tcpListenAddr.String())
        for {
            select {
            case <-s.stop:
                break
            default:
                conn, err := s.tcpListener.AcceptTCP()
                if err != nil {
                    // log this or... somthing
                    continue
                }

                c := NewClient(s, conn)
                c.ConnectionId = utils.NewUUID(conn.RemoteAddr.String())

                s.register <- c
            }
        }
    }()

    go func() {
        // This is for handling UDP data but I ain't got time for that shit right now.
        log.Server.Info("server", "Now listening for UDP at %s", s.config.udpListenAddr.String())

        for {
            select {
            case <-s.stop:
                break
            }
        }
    }()

    return nil
}

func (s *Server) RemoveClient(c *Client) {
    s.unregister <- c
}

func (s *Server) SendUDP(c *Client, p network.Packet) error {
    return nil
}

func (s *Server) cleanup() {
    // This is where we might stop all managers, save to the DB, etc.
}

func (s *Server) Stop() {
    log.Server.Info("server", "Stopping the server...")

    go func() {
        // Spray and pray: send tons of stop messages, so anyone listening is sure to get one..
        for {
            s.stop <- struct{}{}
        }
    }()

    s.cleanup()
}
