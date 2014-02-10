package clientstart

import (
    "database/sql"
	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
    "github.com/tobz/phosphorus/database/models"
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

    // See if the account exists yet.  If so, validate that our password matches.
    tx, err := c.Server().Database().Begin()
    if err != nil {
        log.Server.ClientError(c, "login", "Couldn't start transaction for account validation: %s", err)
        return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
    }
    defer tx.Rollback()

    var account models.Account
    err = tx.SelectOne(&account, "SELECT * FROM account WHERE username = ?", username)
    if err == sql.ErrNoRows {
        log.Server.ClientError(c, "login", "No account found for '%s'!", username)
        return SendLoginDenied(c, constants.LoginErrorAccountNotFound)
    }

    if err != nil {
        log.Server.ClientError(c, "login", "Caught an error when trying to select account: %s", err)
        return SendLoginDenied(c, constants.LoginErrorAccountNotFound)
    }

    // They do, in fact, have an account, so make sure the password matches.
    if !account.PasswordMatches(password) {
        return SendLoginDenied(c, constants.LoginErrorWrongPassword)
    }

    c.SetAccount(&account)

    // At this point, their password matches... so let them in.
	return SendLoginGranted(c)
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

func SendLoginGranted(c interfaces.Client) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.OneWayLoginGranted)

	// Send the client version back, too.  Fucking client is obsessed with it.
	versionDivisor := 100
	if c.ClientVersion() > 199 {
		versionDivisor = 1000
	}

	p.WriteUint8(0x01)
	p.WriteUint8(uint8(int(c.ClientVersion()) / versionDivisor))
	p.WriteUint8(uint8((int(c.ClientVersion()) % versionDivisor) / 10))
	p.WriteUint8(0x00)

	// Write the account name, and then server shortname.
	p.WriteLengthPrefixedString(c.Account().Username())
	p.WriteLengthPrefixedString(c.Server().ShortName())

	// Write ther server ID.
	p.WriteUint8(0x0C)

	// Write the realm color stuff.  The constant definition explains
	// the differences.
	p.WriteUint8(uint8(c.Server().Ruleset().ColorHandling()))
	p.WriteUint8(0x00)

	return c.Send(p)
}
