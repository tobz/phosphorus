package character

import (
	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestCharacterSelect, HandleCharacterSelect)
}

func HandleCharacterSelect(c interfaces.Client, p *network.InboundPacket) error {
	// Skip session ID.
	p.Skip(2)

	// Skip "type."
	p.Skip(2)

	// Skip some random byte.
	p.Skip(1)

	// Grab the character name.
	characterName, err := p.ReadBoundedString(28)
	if err != nil {
		return err
	}

	// Skip these three random bytes.
	p.Skip(3)

	// Skip the user's login name (20 bytes) and client signature (4 bytes)
	p.Skip(24)

	// Skip twenty-four more arbitrary bytes.
	p.Skip(24)

	// Skip the client flag (4 bytes).  This has information about hooked processes and shit, apparently.
	p.Skip(4)

	// If the character name is "noname", they're just asking for a session ID.
	if characterName == "noname" {
		c.Logger().Debug("characterselect", "Character name presented as 'noname' - sending session ID.")

		return SendSessionID(c)
	}

	c.Logger().Debug("characterselect", "Client requested character name: %s.  We don't handle this yet.", characterName)

	return nil
}

func SendSessionID(c interfaces.Client) error {
	c.Logger().Debug("sendsessionid", "Sending session ID as %d...", c.SessionID())

	p := network.NewOutboundPacket(constants.PacketTCP, constants.ServerOnlySessionID)
	p.WriteUInt16(c.SessionID())

	return c.Send(p)
}
