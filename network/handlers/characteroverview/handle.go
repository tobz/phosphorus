package characteroverview

import (
	"fmt"
	"strings"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/database/models"
	"github.com/tobz/phosphorus/helpers"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestCharacterOverview, HandleCharacterOverview)
}

func HandleCharacterOverview(c interfaces.Client, p *network.InboundPacket) error {
	// Our client is at the character screen now, so update their state.
	c.SetClientState(constants.ClientStateCharacterScreen)

	// Read their account name.  We technically already have it, but this one has a suffix
	// that indicates which realm the client wants to play.
	accountName, err := p.ReadBoundedString(24)
	if err != nil {
		return err
	}

	c.Logger().Debug("characteroverview", "Account name reported as: %s", accountName)

	// If the account name is prefixed with "-X", it means the client wants us to tell them whether
	// or not they can pick from any realm.  If they're already assigned a realm, and this isn't
	// a server which allows a player to choose any realm (like PvP) then we tell them which realm
	// they belong to.  If they have no realm assigned, either from being a new account or having
	// zero characters, we let them pick one.
	if strings.HasSuffix(accountName, "-X") {
		c.Logger().Debug("characteroverview", "Client reported account name with -X suffix: sending realm.")

		if c.Server().Ruleset().CanChooseAnyRealm(c) {
			return SendRealm(c, constants.ClientRealmNone)
		} else {
			switch constants.ClientRealm(c.Account().Realm) {
			case constants.ClientRealmAlbion:
				return SendRealm(c, constants.ClientRealmAlbion)
			case constants.ClientRealmMidgard:
				return SendRealm(c, constants.ClientRealmMidgard)
			case constants.ClientRealmHibernia:
				return SendRealm(c, constants.ClientRealmHibernia)
			default:
				return SendRealm(c, constants.ClientRealmNone)
			}
		}
	} else {
		c.Logger().Debug("characteroverview", "Client reported account name without -X suffix: figuring out proper realm.")

		// The client picked a realm from the realm selection screen, and they're now telling us
		// what their selection was.  Make sure it's a valid selection and then assign them and
		// send them the character list overview.
		startingRealm := constants.ClientRealm(c.Account().Realm)

		chosenRealm := startingRealm
		chosenRealmName := ""

		c.Logger().Debug("characteroverview", "Client's current realm association: %d", c.Account().Realm)

		// If they're set to "none" - it means they have a new account.  Their choice here will either stick or not,
		// depending on the ruleset.  If they aren't set to none, but this is an open realm ruleset (PvP), we'll
		// get their choice and use that for this session.
		if chosenRealm == constants.ClientRealmNone || c.Server().Ruleset().CanChooseAnyRealm(c) {
			if strings.HasSuffix(accountName, "-S") {
				chosenRealm = constants.ClientRealmAlbion
				chosenRealmName = "albion"
			} else if strings.HasSuffix(accountName, "-N") {
				chosenRealm = constants.ClientRealmMidgard
				chosenRealmName = "midgard"
			} else if strings.HasSuffix(accountName, "-H") {
				chosenRealm = constants.ClientRealmHibernia
				chosenRealmName = "hibernia"
			} else {
				// Unknown suffix.  Bye bye.
				return fmt.Errorf("unknown suffix or no suffix supplied for character overview request.  Disconnecting.")
			}

			// If they started out with no realm association, and this isn't an open realm ruleset, set their account
			// to the realm they selected.
			if startingRealm == constants.ClientRealmNone && !c.Server().Ruleset().CanChooseAnyRealm(c) {
				c.Logger().Debug("characteroverview", "Ruleset prohibits multiple realms: setting account realm as '%s'", chosenRealmName)

				c.Account().Realm = uint8(chosenRealm)
				if _, err = c.Server().Database().Update(c.Account()); err != nil {
					return fmt.Errorf("caught an error while saving realm selection: %s", err)
				}
			}
		}

		return SendCharacterOverview(c, chosenRealm)
	}
}

func SendRealm(c interfaces.Client, realm constants.ClientRealm) error {
	c.Logger().Debug("sendrealm", "Sending realm as %X...", realm)

	p := network.NewOutboundPacket(constants.PacketTCP, constants.ServerOnlyRealm)
	p.WriteUInt8(uint8(realm))
	return c.Send(p)
}

func SendCharacterOverview(c interfaces.Client, realm constants.ClientRealm) error {
	c.Logger().Debug("sendcharacteroverview", "Sending character overview for realm %d...", realm)

	p := network.NewOutboundPacket(constants.PacketTCP, constants.ResponseCharacterOverview)

	// Write their account name.
	p.WriteBoundedString(c.Account().Username, 28)

	// Get all the characters for this account so we can find the ones for the given realm.
	characters, err := helpers.GetCharactersForClient(c)
	if err != nil {
		return err
	}

	// If we have no characters, just send an empty packet.
	if len(characters) == 0 {
		p.WriteRepeated(0x00, 1950)
		return c.Send(p)
	}

	// Build a map of account slots for characters so we can properly display characters
	// in the slots they should be in.
	charactersBySlot := make(map[int]*models.Character)

	for _, character := range characters {
		charactersBySlot[int(character.AccountSlot)] = character
	}

	// Now iterate through the map, going in order, so we can determine an empty slot from
	// an occupied slot.
	startingSlot := int(realm) * 100

	for slot := startingSlot; slot < (startingSlot + 10); slot++ {
		character, ok := charactersBySlot[slot]
		if !ok {
			// No character in the slot, so we just send a bunch of zeros.
			p.WriteRepeated(0x00, 188)
			continue
		}

		// Write the character name.
		p.WriteBoundedString(character.FirstName, 24)

		// Random byte.
		p.WriteUInt8(0x01)

		// Write all the facial attributes.
		p.WriteUInt8(character.EyeSize)
		p.WriteUInt8(character.LipSize)
		p.WriteUInt8(character.EyeColor)
		p.WriteUInt8(character.HairColor)
		p.WriteUInt8(character.FaceType)
		p.WriteUInt8(character.HairStyle)

		// These control armor extensions, but we don't have that yet, so...
		p.WriteUInt8(0x00)
		p.WriteUInt8(0x00)

		// More facial stuff.  Why is this broken up? :/
		p.WriteUInt8(character.MoodType)

		// Random filler.
		p.WriteRepeated(0x00, 13)

		// This is supposed to be where the character is located. No support yet, so this is an
		// empty string for now.
		p.WriteBoundedString("", 24)

		// Class and race name.
		p.WriteBoundedString(helpers.GetClassName(character.Class, true), 24)
		p.WriteBoundedString(helpers.GetRaceName(character.Race, true), 24)

		// Some more raw values - realm, class, level, etc.
		p.WriteUInt8(character.Level)
		p.WriteUInt8(character.Class)
		p.WriteUInt8(character.Realm)
		p.WriteUInt8(uint8(((character.Race&0x10)<<2)+(character.Race&0x0F)) | uint8(character.Gender<<4))

		// Time for some models and shit.  Region is mixed in here, too.
		p.WriteHUInt16(character.CurrentModel)

		// Region is 0 for now. Write another 0 since we don't have region data and the expansion values sorted yet.
		p.WriteUInt8(0x00)
		p.WriteUInt8(0x00)

		// Random integer.  Supposed to be internal ID?  Dunno.  Fuck it.
		p.WriteUInt32(0x00)

		// Character stats.  Kind of weird since these can be greater (2^8)-1 but whatever.
		p.WriteUInt8(uint8(character.Strength))
		p.WriteUInt8(uint8(character.Dexterity))
		p.WriteUInt8(uint8(character.Constitution))
		p.WriteUInt8(uint8(character.Quickness))
		p.WriteUInt8(uint8(character.Intelligence))
		p.WriteUInt8(uint8(character.Piety))
		p.WriteUInt8(uint8(character.Empathy))
		p.WriteUInt8(uint8(character.Charisma))

		// Here comes the item values.  We definitely don't have an item system yet so these are all going to be zero.
		p.WriteHUInt16(0x00) // Helmet model.
		p.WriteHUInt16(0x00) // Gloves model.
		p.WriteHUInt16(0x00) // Boots model.
		p.WriteHUInt16(0x00) // RH weapon color.
		p.WriteHUInt16(0x00) // Torso model.
		p.WriteHUInt16(0x00) // Cloak model.
		p.WriteHUInt16(0x00) // Legs model.
		p.WriteHUInt16(0x00) // Arms model.

		p.WriteHUInt16(0x00) // Helmet color.
		p.WriteHUInt16(0x00) // Gloves color.
		p.WriteHUInt16(0x00) // Boots color.
		p.WriteHUInt16(0x00) // LH weapon color.
		p.WriteHUInt16(0x00) // Torso color.
		p.WriteHUInt16(0x00) // Cloak color.
		p.WriteHUInt16(0x00) // Legs color.
		p.WriteHUInt16(0x00) // Arms model.

		p.WriteHUInt16(0x00) // RH weapon model.
		p.WriteHUInt16(0x00) // LH weapon model.
		p.WriteHUInt16(0x00) // Two-hand weapon model.
		p.WriteHUInt16(0x00) // Ranged weapon model.

		// These control which weapons are currently equipped.  0xFF for both means "no weapons equipped."
		p.WriteUInt8(0xFF)
		p.WriteUInt8(0xFF)

		// No regions / expansion data yet, so another 0 here.
		p.WriteUInt8(0x00)

		// Write the character's constitution again... no idea why.
		p.WriteUInt8(uint8(character.Constitution))

		// Four filler bytes.
		p.WriteRepeated(0x00, 4)
	}

	// Random filler down here.
	p.WriteRepeated(0x00, 90)

	return c.Send(p)
}
