package interfaces

type Client interface {
    SetAccount(Account)
	Account() Account

    ConnectionId() uint32

    Server() Server

	Send(Packet) error
}
