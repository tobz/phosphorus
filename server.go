package phosphorus

import "net"
import "sync"
import "github.com/tobz/phosphorus/network"

type Server struct {
	config *ServerConfig

	tcpListener *net.TCPListener
	udpListener *net.UDPConn
	coordinator *Coordinator

	clients    map[string]*Client
	clientLock *sync.RWMutex
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		config:      config,
		coordinator: NewCoordinator(),
		tcpListener: nil,
		udpListener: nil,
		clients:     make(map[string]*Client, 128),
		clientLock:  &sync.RWMutex{},
	}
}

func (s *Server) Start() error {
	// Do a bunch of random shit here.

	// Now start listening.
	tl, err := net.ListenTCP(s.config.tcpListenAddr.Network(), s.config.tcpListenAddr)
	if err != nil {
		s.Cleanup()
		return err
	}

	s.tcpListener = tl

	ul, err := net.ListenUDP(s.config.udpListenAddr.Network(), s.config.udpListenAddr)
	if err != nil {
		s.Cleanup()
		return err
	}

	s.udpListener = ul

	// Now run our workers that pull things off the wire and accept connections.
	go func() {
		// Register ourselves with the coordinator.
		stop := s.coordinator.Register()

		for {
			select {
			case <-stop:
				break
			default:
				c, err := s.tcpListener.AcceptTCP()
				if err != nil {
					// Log this or... something.

					continue
				}

				go s.handleNewConnection(c)
			}
		}
	}()

	go func() {
		// This is for handling UDP data but I ain't got
		// time for that shit right now.
		stop := s.coordinator.Register()

		for {
			select {
			case <-stop:
				break
			}
		}
	}()

	return nil
}

func (s *Server) handleNewConnection(connection *net.TCPConn) {
	newClient := NewClient(s, connection)

	s.clientLock.Lock()
	s.clients[newClient.GetUniqueIdentifier()] = newClient
	s.clientLock.Unlock()

	newClient.Start()
}

func (s *Server) SendUDP(client *Client, packet interfaces.Packet) error {
	return nil
}

func (s *Server) Cleanup() {
	// This is where we might stop all managers, save to the DB, etc.
}

func (s *Server) Stop() {
	// Stop all the worker routines and what not.
	s.coordinator.Ping()
	s.coordinator.Stop()

	// Now clean up after ourselves.
	s.Cleanup()
}
