package phosphorus

import "net"
import "sync"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/utils"

type Server struct {
	config *ServerConfig

	tcpListener *net.TCPListener
	udpListener *net.UDPConn
	coordinator *utils.Coordinator

	clients    map[uint32]*Client
	clientLock *sync.RWMutex
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		config:      config,
		coordinator: utils.NewCoordinator(),
		tcpListener: nil,
		udpListener: nil,
		clients:     make(map[uint32]*Client, 128),
		clientLock:  &sync.RWMutex{},
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
		// Register ourselves with the coordinator.
		stop := s.coordinator.Register()

		log.Server.Info("server", "Now listening for TCP connections at %s", s.config.tcpListenAddr.String())

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

		log.Server.Info("server", "Now listening for UDP at %s", s.config.udpListenAddr.String())

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
	log.Server.Info("server", "Accepting new connection from client %s", connection.RemoteAddr().String())

	newClient := NewClient(s, connection)

	s.AddClient(newClient)

	newClient.Start()
}

func (s *Server) AddClient(client *Client) {
	s.clientLock.Lock()
	s.clients[client.GetUniqueIdentifier()] = client
	s.clientLock.Unlock()
}

func (s *Server) RemoveClient(client *Client) {
	s.clientLock.Lock()
	delete(s.clients, client.GetUniqueIdentifier())
	s.clientLock.Unlock()
}

func (s *Server) SendUDP(client *Client, packet interfaces.Packet) error {
	return nil
}

func (s *Server) Cleanup() {
	// This is where we might stop all managers, save to the DB, etc.
}

func (s *Server) Stop() {
	log.Server.Info("server", "Stopping the server...")

	// Stop all the worker routines and what not.
	s.coordinator.Ping()
	s.coordinator.Stop()

	// Now clean up after ourselves.
	s.Cleanup()
}
