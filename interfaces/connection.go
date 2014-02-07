package interfaces

import "net"

type LimitedConnection interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
