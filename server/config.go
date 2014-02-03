package phosphorus

import (
	"net"
)

type ServerConfig struct {
	TcpListenAddress *net.TCPAddr
	UdpListenAddress *net.UDPAddr

    InvalidWords []string
}
