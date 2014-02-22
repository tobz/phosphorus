package character

import (
	"fmt"
	"regexp"
	"strings"
    "database/sql"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/database/models"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestCharacterCreate, HandleCharacterCreate)
}

func HandleCharacterCreate(c interfaces.Client, p *network.InboundPacket) error {
	// Get their account name.
	accountName, err := p.ReadBoundedString(24)
	if err != nil {
		return err
	}

	// Make sure it matches their existing account name.  Cheap insurance.
	if !strings.HasPrefix(accountName, c.Account().Username) {
		return fmt.Errorf("client supplied erroneous account name: given account name '%s' expected to start with '%s'", accountName, c.Account().Username)
	}

	// Skip four random bytes.  Dunno what they're for.
	p.Skip(4)

	// Loop through all the possible slots.  There will either be character info or filler for the
	// given slot.  The differences are explained below.
    characterRealm := constants.ClientRealmNone

	for slot := 0; slot < 10; slot++ {
		// Try and get a character name.
		characterName, err := p.ReadBoundedString(24)
		if err != nil {
			return err
		}

		if characterName == "" {
			// If the name is empty, it means the client thinks there shouldn't be a character in
			// this slot.  We need to check the database to see if there actually is a character
			// in this slot, and if so, delete them.  We also have to skip all the filler data.
			deleteCharacterIfExists(c, accountName, characterName, slot)

			p.Skip(164)
		} else {
			// We have a character!  Let's verify their name is valid and all of that.
			matches, err := regexp.MatchString("^[A-Z][a-z]+$", characterName)
			if err != nil {
				return err
			}

			// If the character name doesn't match, it means the client is doing something funky
			// that shouldn't be possible, so we should disconnect these fools.
			if !matches {
				return fmt.Errorf("client supplied invalid name during character creation: %s", characterName)
			}

			// Now see if this is an existing character.  Make sure it belongs to us, too.
			tx, err := c.Server().Database().Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			character := &models.Character{}
			err = tx.SelectOne(character, "SELECT * FROM character WHERE first_name = ?", characterName)
			if err == nil {
				// We found a character.  Make sure they belong to us before trying to customize.
				if character.AccountID != c.Account().AccountID {
					return fmt.Errorf("client tried to run customization for a character on another account: %s", characterName)
				}

				// Try and do the customization.  Close our transaction first since we don't need
				// it anymore.  There's no risk of multi-client collision now.
				tx.Rollback()

				handleCharacterCustomization(c, p, accountName, characterName, slot)
			} else if err != nil && err != sql.ErrNoRows {
				// We got a legitimate error.  Wah.
				return err
			}

			// We didn't find another existing character.  Proceed with trying to create the character.  Pull
            // out the realm this character is so we can send the character overview at the end.
            characterRealm = handleCharacterCreate(c, p, accountName, characterName, slot)
		}
	}

	return SendCharacterOverview(c, characterRealm)
}

func deleteCharacterIfExists(c interfaces.Client, accountName, characterName string, slot int) {
}

func handleCharacterCustomization(c interfaces.Client, p *network.InboundPacket, accountName, characterName string, slot int) {
}

func handleCharacterCreate(c interfaces.Client, p *network.InboundPacket, accountName, characterName string, slot int) constants.ClientRealm {
    return constants.ClientRealmNone
}
