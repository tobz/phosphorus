package server

import (
	"net"

	"github.com/rcrowley/go-metrics"

	"github.com/tobz/phosphorus/database"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/rulesets"
	"github.com/tobz/phosphorus/statistics"
	"github.com/tobz/phosphorus/utils"
	"github.com/tobz/phosphorus/world"
)

type Server struct {
	config   *Config
	database interfaces.Database

	tcpListener *net.TCPListener
	udpListener *net.UDPConn

	clients map[uint32]*Client

	register   chan *Client
	unregister chan *Client
	stop       chan struct{}

	ruleset interfaces.Ruleset
	world   interfaces.World

	connectionCount metrics.Counter
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,

		clients: make(map[uint32]*Client, 128),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan struct{}),

		connectionCount: metrics.GetOrRegisterCounter("server.connectionCount", statistics.Registry),
	}
}

func (s *Server) Start() error {
	log.Server.Info("server", "Starting the server...")

	// This is where managers and all that shit will be instanciated.

	// Set up our statistics if they're configured.
	statisticsType, err := s.config.GetAsString("statistics/provider")
	if err != nil {
		// Just warn the user that statistics won't be running.  It's not a reason to *not* run the server.
		log.Server.Warn("server", "Statistics configuration missing: server will not record statistics.  Update statistics configuration and restart Phosphorus to enable.")
	} else {
		switch statisticsType {
		case "influxdb":
			err = statistics.ConfigureInfluxDB(s.config)
			if err != nil {
				return err
			}
		default:
			log.Server.Warn("server", "Unknown statistics provider '%s': server will not record statistics.  Update statistics configuration and restart Phosphorus to enable.")
		}
	}

	// Get our database connection poppin'.
	databaseType, err := s.config.GetAsString("database/type")
	if err != nil {
		return err
	}

	databaseDsn, err := s.config.GetAsString("database/dsn")
	if err != nil {
		return err
	}

	database, err := database.NewDatabaseConnection(databaseType, databaseDsn)
	if err != nil {
		return err
	}

	s.database = database

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

	// Load our world manager.
	world, err := world.NewWorldMgr(s.config)
	if err != nil {
		return err
	}

	err = world.Start()
	if err != nil {
		return err
	}

	s.world = world

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
				s.clients[c.ConnectionID()] = c

				s.connectionCount.Inc(1)

				c.Start()
			case c := <-s.unregister:
				log.Server.Info("server", "Closing connection for client %s", c.connection.RemoteAddr().String())
				delete(s.clients, c.ConnectionID())
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

func (s *Server) Database() interfaces.Database {
	return s.database
}

func (s *Server) World() interfaces.World {
	return s.world
}

func (s *Server) ShortName() string {
	shortName, err := s.config.GetAsString("server/shortName")
	if err != nil {
		return "noname"
	}

	return shortName
}
