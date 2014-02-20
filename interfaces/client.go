package interfaces

import "time"
import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/database/models"

type Client interface {
	Connection() LimitedConnection

	Logger() Logger

	SetAccount(*models.Account)
	Account() *models.Account

	SetClientVersion(constants.ClientVersion)
	ClientVersion() constants.ClientVersion
	SetClientState(constants.ClientState)
	ClientState() constants.ClientState

	ConnectionID() uint32

	SetSessionID(uint16)
	SessionID() uint16

	LastPingTime() time.Time
	MarkPingTime()

	Server() Server

	Stop()
	Send(Packet) error
}
