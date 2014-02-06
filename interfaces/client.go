package interfaces

import "github.com/tobz/phosphorus/constants"

type Client interface {
    Connection() LimitedConnection

    SetAccount(Account)
	Account() Account

    SetClientVersion(constants.ClientVersion)
    ClientVersion() constants.ClientVersion

    ConnectionId() uint32

    Server() Server

	Send(Packet) error
}
