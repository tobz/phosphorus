package phosphorus

import (
	"net"
)

type ServerConfig struct {
	tcpListenAddr *net.TCPAddr
	udpListenAddr *net.UDPAddr
}
