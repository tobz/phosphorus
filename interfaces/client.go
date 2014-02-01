package interfaces

import "net"

type Client interface {
	SetAccount(account Account)
	Account() Account

	RemoteAddr() net.Addr

	Send(packet Packet) error
}
