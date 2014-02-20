package clientstart

import (
	"database/sql"
	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/database/models"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
	"time"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestLogin, HandleLoginRequest)
}

func HandleLoginRequest(c interfaces.Client, p *network.InboundPacket) error {
	// Not sure what this is. *shrug*
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

	c.Logger().Debug("loginrequest", "Attempting login of '%s'...", username)

	// See if the account exists yet.  If so, validate that our password matches.
	tx, err := c.Server().Database().Begin()
	if err != nil {
		c.Logger().Error("login", "Couldn't start transaction for account validation: %s", err)
		return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
	}
	defer tx.Rollback()

	account := &models.Account{}
	err = tx.SelectOne(account, "SELECT * FROM accounts WHERE username = ?", username)
	if err == sql.ErrNoRows {
		autocreate, err := c.Server().Config().GetAsBoolean("server/autocreateAccounts")
		if err != nil {
			log.Server.Debug("login", "Unable to load 'server/autocreateAccounts' - automatic account creation will fail until present.")
			return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
		}

		if !autocreate {
			c.Logger().Error("login", "No account found for '%s'!", username)
			return SendLoginDenied(c, constants.LoginErrorAccountNotFound)
		}

		c.Logger().Error("login", "No account found for '%s': creating...", username)

		account, err = models.NewAccount(username, password)
		if err != nil {
			c.Logger().Error("login", "Caught an error while creating a new account: %s", err)
			return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
		}

		// Since we got the account, we know we're going to login in so set the last login time now.
		account.LastLogin = time.Now()

		// We're all set here, so save the account and close our transaction.
		err = tx.Insert(account)
		if err != nil {
			c.Logger().Debug("login", "Caught an error while saving a new account: %s", err)
			return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
		}

		err = tx.Commit()
		if err != nil {
			c.Logger().Debug("login", "Caught an error trying to commit our transaction: %s", err)
			return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
		}
	} else {
		if err != nil {
			c.Logger().Error("login", "Caught an error when trying to select account: %s", err)
			return SendLoginDenied(c, constants.LoginErrorAccountNotFound)
		}

		// They do, in fact, have an account, so make sure the password matches.
		matches, err := account.PasswordMatches(password)
		if err != nil {
			c.Logger().Error("login", "Caught an error while trying to match passwords for an account: %s", err)
			return SendLoginDenied(c, constants.LoginErrorAuthorizationServerUnavailable)
		}

		if !matches {
			c.Logger().Error("loginrequest", "User '%s' entered wrong password!", username)
			return SendLoginDenied(c, constants.LoginErrorWrongPassword)
		}
	}

	c.Logger().Debug("loginrequest", "User authenticated successfully!  Adding to world...")

	c.SetAccount(account)

	// Try and add them to the world.  Even though they don't have a character selected,
	// they still need a session ID and what not for some packets that happen before
	// actually playing.
	err = c.Server().World().AddClient(c)
	if err != nil {
		c.Logger().Debug("loginrequest", "Failed to assign session ID to client.")

		SendLoginDenied(c, constants.LoginErrorTooManyPlayersLoggedIn)
		return err
	}

	// At this point, their password matches... so let them in.
	return SendLoginGranted(c)
}

func SendLoginDenied(c interfaces.Client, reason constants.LoginError) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.ServerOnlyLoginDenied)

	// Write the reason we're denying them.
	p.WriteUInt8(uint8(reason))

	// Send the client version back, too.  Fucking client is obsessed with it.
	versionDivisor := 100
	if c.ClientVersion() > 199 {
		versionDivisor = 1000
	}

	p.WriteUInt8(0x01)
	p.WriteUInt8(uint8(int(c.ClientVersion()) / versionDivisor))
	p.WriteUInt8(uint8((int(c.ClientVersion()) % versionDivisor) / 10))
	p.WriteUInt8(0x00)

	return c.Send(p)
}

func SendLoginGranted(c interfaces.Client) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.ServerOnlyLoginGranted)

	// Write the account name, and then server shortname.
	p.WriteLengthPrefixedString(c.Account().Username)
	p.WriteLengthPrefixedString(c.Server().ShortName())

	// Write ther server ID.
	p.WriteUInt8(0x0C)

	// Write the realm color stuff.  The constant definition explains
	// the differences.
	p.WriteUInt8(uint8(c.Server().Ruleset().ColorHandling()))
	p.WriteUInt8(0x00)

	// Write the "server index." Not sure what the hell this is, so let's just fake it?
	p.WriteUInt8(0x00)

	return c.Send(p)
}
