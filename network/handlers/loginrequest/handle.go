package clientstart

import (
	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
    "github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestLogin, HandleLoginRequest)
}

func HandleLoginRequest(c interfaces.Client, p *network.InboundPacket) error {
    p.Skip(2)

    // Skip the client version.
    p.Skip(3)

    // Read the password.
    password, err := p.ReadBoundedString(20)
    if err != nil {
        return err
    }

    // Skip some of this cruft in here.
    p.Skip(50)

    // Read the username.
    username, err := p.ReadBoundedString(20)
    if err != nil {
        return err
    }

    log.Server.ClientDebug(c, "login", "username: %s(%d), password: %s(%d)", username, len(username), password, len(password))

    // From here, we have many options.  We can respond back with any of the multiple
    // login errors - game closed, already logged in, account mid-logout, etc.

    return SendLoginDenied(c, constants.LoginErrorGameCurrentlyClosed)
}

func SendLoginDenied(c interfaces.Client, reason constants.LoginError) error {
    p := network.NewOutboundPacket(constants.PacketTCP, constants.OneWayLoginDenied)

    // Write the reason we're denying them.
    p.WriteUint8(uint8(reason))

    // Send the client version back, too.  Fucking client is obsessed with it.
    versionDivisor := 100
    if c.ClientVersion() > 199 {
        versionDivisor = 1000
    }

    p.WriteUint8(0x01)
    p.WriteUint8(uint8(int(c.ClientVersion()) / versionDivisor))
    p.WriteUint8(uint8((int(c.ClientVersion()) % versionDivisor) / 10))
    p.WriteUint8(0x00)

    return c.Send(p)
}

