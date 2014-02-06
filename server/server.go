package server

import (
	"net"

	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/utils"
    "github.com/tobz/phosphorus/rulesets"
)

type Server struct {
	config *Config

	tcpListener *net.TCPListener
	udpListener *net.UDPConn

	clients map[uint32]*Client

	register   chan *Client
	unregister chan *Client
	stop       chan struct{}

    ruleset interfaces.Ruleset
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,

		tcpListener: nil,
		udpListener: nil,

		clients: make(map[uint32]*Client, 128),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan struct{}),

        ruleset: nil,
	}
}

func (s *Server) Start() error {
	log.Server.Info("server", "Starting the server...")

    // This is where managers and database connections and all that shit will be instanciated.

	// Load the ruleset we should be using.
    rulesetName, err := s.config.GetAsString("server/ruleset")
    if err != nil {
        return err
    }

    ruleset, err := rulesets.GetRuleset(rulesetName, s)
    if err != nil {
        return err
    }

    s.ruleset = ruleset

    log.Server.Info("server", "Ruleset configured for '%s'", rulesetName)

	// Now start listening.
	tcpListenAddr, err := s.config.GetAsString("server/tcpListen")
	if err != nil {
		s.Stop()
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpListenAddr)
	if err != nil {
		s.Stop()
		return err
	}

	tl, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	if err != nil {
		s.Stop()
		return err
	}

	s.tcpListener = tl

	udpListenAddr, err := s.config.GetAsString("server/udpListen")
	if err != nil {
		s.Stop()
		return err
	}

	udpAddr, err := net.ResolveUDPAddr("udp", udpListenAddr)
	if err != nil {
		s.Stop()
		return err
	}

	ul, err := net.ListenUDP(udpAddr.Network(), udpAddr)
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
				log.Server.Info("server", "Accepting new connection from client %s", c.connection.RemoteAddr().String())
				s.clients[c.ConnectionId()] = c

				c.Start()
			case c := <-s.unregister:
				log.Server.Info("server", "Closing connection for client %s", c.connection.RemoteAddr().String())
				delete(s.clients, c.ConnectionId())
			}
		}
	}()

	go func() {
		log.Server.Info("server", "Now listening for TCP connections at %s", tcpAddr.String())
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

				s.register <- NewClient(s, conn, utils.NewUUID(conn.RemoteAddr().String()))
			}
		}
	}()

	go func() {
		// This is for handling UDP data but I ain't got time for that shit right now.
		log.Server.Info("server", "Now listening for UDP at %s", udpAddr.String())

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

func (s *Server) SendUDP(c *Client, p interfaces.Packet) error {
	return nil
}

func (s *Server) cleanup() {
	// This is where we might stop all managers, save to the DB, etc.
}

func (s *Server) Stop() {
	log.Server.Info("server", "Stopping the server...")

	go func() {
		// Spray and pray: send tons of stop messages, so anyone listening is sure to get one.
		for {
			s.stop <- struct{}{}
		}
	}()

	s.cleanup()
}

// Methods to satisfy interfaces.Server
func (s *Server) Config() interfaces.Config {
	return s.config
}

func (s *Server) Ruleset() interfaces.Ruleset {
    return s.ruleset
}

func (s *Server) ShortName() string {
    shortName, err := s.config.GetAsString("server/shortName")
    if err != nil {
        return "noname"
    }

    return shortName
}
