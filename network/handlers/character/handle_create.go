package character

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
    "log"

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
	characterRealm := constants.ClientRealm(c.Account().Realm)

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
            log.Printf("Character found: %s", characterName)

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
			err = tx.SelectOne(character, "SELECT * FROM characters WHERE first_name = ?", characterName)
            switch {
            case err == nil:
				// We found a character.  Make sure they belong to us before trying to customize.
				if character.AccountID != c.Account().AccountID {
					return fmt.Errorf("client tried to run customization for a character on another account: %s", characterName)
				}

				// Do the customization.
				handleCharacterCustomization(c, p, accountName, characterName, slot)
            case err == sql.ErrNoRows:
                // We didn't find another existing character.  Proceed with trying to create the character.  Pull
                // out the realm this character is so we can send the character overview at the end.
                character, err = handleCharacterCreate(c, p, accountName, characterName, slot)
                if err != nil {
                    return err
                }

                err = tx.Insert(character)
                if err != nil {
                    return err
                }

                // Close our transaction.
                err = tx.Commit()
                if err != nil {
                    // Can we send something here to drop them back to the login error screen?
                    return err
			    }
            default:
				// We got a legitimate error.  Wah.
				return err
			}
		}
	}

	return SendCharacterOverview(c, characterRealm)
}

func deleteCharacterIfExists(c interfaces.Client, accountName, characterName string, slot int) {
}

func handleCharacterCustomization(c interfaces.Client, p *network.InboundPacket, accountName, characterName string, slot int) {
    p.Skip(164)
}

func handleCharacterCreate(c interfaces.Client, p *network.InboundPacket, accountName, characterName string, slot int) (*models.Character, error) {
	var err error

	customization, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	// Create our character.
	character := &models.Character{}
	character.AccountID = c.Account().AccountID
	character.GuildID = 0
	character.FirstName = characterName
	character.Level = 1

	// Get all of the customization details, if any.
	if customization != 0 {
		character.EyeSize, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		character.LipSize, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		character.EyeColor, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		character.HairColor, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		character.FaceType, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		character.HairStyle, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		p.Skip(3)

		character.MoodType, err = p.ReadUInt8()
		if err != nil {
			return nil, err
		}

		p.Skip(13)
	} else {
		// Not doing any customizing?  Just do some skippin' then.
		p.Skip(23)
	}

	// Skip the location string, class name and race name.
	p.Skip(72)

	// Skip the level byte because we're just going to start them at level 1.
	p.Skip(1)

	// Get their class and realm.
	character.Class, err = p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Realm, err = p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	// Get their race, gender
	startRaceGender, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Race = (startRaceGender & 0x0F) + ((startRaceGender & 0x40) >> 2)
	character.Gender = ((startRaceGender >> 4) & 0x01)
	startInShroudedIsles := ((startRaceGender >> 7) != 0)

	// Get their model and region.
	character.BaseModel, err = p.ReadHUInt16()
	if err != nil {
		return nil, err
	}

	character.CurrentModel = character.BaseModel

	character.Region, err = p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	// Now read in their stats.
	strength, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Strength = uint32(strength)

	dexterity, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Dexterity = uint32(dexterity)

	constitution, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Constitution = uint32(constitution)

	quickness, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Quickness = uint32(quickness)

	intelligence, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Intelligence = uint32(intelligence)

	piety, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Piety = uint32(piety)

	empathy, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Empathy = uint32(empathy)

	charisma, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	character.Charisma = uint32(charisma)

	// Skip some random cruft.
	p.Skip(48)

	// Make sure their base class is set if necessary.
	adjustBaseClassIfNecessary(character)

	// Make sure the character so far is valid.  Valid realm, class, race, stats, etc.
	if err = verifyValidCharacter(character); err != nil {
		return nil, err
	}

	// Set some particulars for the character.
	character.Created = time.Now()
	character.AccountSlot = uint32(slot) + uint32(character.Realm * 100)
	character.Endurance = 100
	character.MaxEndurance = 100
	character.MaxSpeed = constants.BaseCharacterSpeed

	// Set their starting point.
	setStartingLocation(character, startInShroudedIsles)

	// Set their starting guild.
	setStartingGuild(character)

	// Now we should be able to save our character so pass it back to the caller.
	return character, nil
}

func adjustBaseClassIfNecessary(character *models.Character) {
}

func setStartingLocation(character *models.Character, startInSi bool) {
}

func setStartingGuild(character *models.Character) {
}

func verifyValidCharacter(character *models.Character) error {
	return nil
}
