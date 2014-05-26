package character

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/database/models"
	"github.com/tobz/phosphorus/helpers"
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

				// Rollback since we're not trying to lock the character name for use.
				tx.Rollback()

				// Do the customization.
				err = handleCharacterCustomization(c, p, character)
				if err != nil {
					return err
				}
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

func handleCharacterCustomization(c interfaces.Client, p *network.InboundPacket, character *models.Character) error {
	customizationStep, err := p.ReadUInt8()
	if err != nil {
		return err
	}

	if customizationStep >= 1 && customizationStep <= 3 {
		statsChanged := false

		character.EyeSize, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		character.LipSize, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		character.EyeColor, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		character.HairColor, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		character.FaceType, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		character.HairStyle, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		p.Skip(3)

		character.MoodType, err = p.ReadUInt8()
		if err != nil {
			return err
		}

		// Skip location, race, class, level, etc etc.
		p.Skip(89)

		newModel, err := p.ReadHUInt16()
		if err != nil {
			return err
		}

		if customizationStep != 3 {
			// Skip region ID and character internal ID.
			p.Skip(6)

			// Read in all of the stats the client sent.  These are: strength, dexterity, constitution,
			// quickness, intelligence, piety, empathy and charisma, in that order.
			p.Skip(8)

			// Skip over armor and weapon values, and the dangling twice-sent constitution value.
			p.Skip(48)

			// Here is where we would see if any stats changed, and if so, make sure the totals
			// were valid, etc, and then apply them.  The DOL code for this s apparently jacked
			// up, so it's a bad reference and this is kind of a corner case in terms of normal
			// players, so we'll figure it out later.  #TODO
		} else {
			// Skip over a bunch of junk if we're not modifying any stats.
			p.Skip(62)
		}

		switch customizationStep {
		case 2:
			// Make sure they have a valid size for their model.
			if ((newModel >> 11) & 3) == 0 {
				return fmt.Errorf("Character size for '%s' set as 0 during customization step!", character.FirstName)
			}

			if character.BaseModel != newModel {
				character.BaseModel = newModel
				character.CurrentModel = newModel
			}

			// We're done here.  Lock down the customization step for this character.
			character.CustomizationStep = 2
		case 3:
			// This apparently re-enables the config button for the player.  We're just following DOL
			// at this point.
			character.CustomizationStep = 3
		}

		if statsChanged || customizationStep != 1 {
			_, err := c.Server().Database().Update(character)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func handleCharacterCreate(c interfaces.Client, p *network.InboundPacket, accountName, characterName string, slot int) (*models.Character, error) {
	var err error

	customizationMode, err := p.ReadUInt8()
	if err != nil {
		return nil, err
	}

	// Create our character.
	character := &models.Character{}
	character.AccountID = c.Account().AccountID
	character.GuildID = 0
	character.FirstName = characterName
	character.Level = 1
	character.CustomizationStep = 1

	// Get all of the customization details, if any.
	if customizationMode == 1 {
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

	// Skip 5 - 2nd byte of region and unknown int.
	p.Skip(5)

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
	character.AccountSlot = uint32(slot) + uint32(character.Realm*100)
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
	characterRealm := constants.ClientRealm(character.Realm)
	characterClass := constants.CharacterClass(character.Class)
	characterRace := constants.CharacterRace(character.Race)
	classFull := helpers.ClassToString(characterClass)
	realmFull := helpers.RealmToString(characterRealm)
	raceFull := helpers.RaceToString(characterRace)

	if characterRealm < constants.ClientRealmMinimum || characterRealm > constants.ClientRealmMaximum {
		return fmt.Errorf("New character realm is out-of-bounds: %d", characterRealm)
	}

	if character.Level != 1 {
		return fmt.Errorf("New character level is %d, should be 1", character.Level)
	}

	classesForRealm, ok := constants.CharacterStartingClasses[characterRealm]
	if !ok {
		return fmt.Errorf("Unable to get starting classes for realm %s (%d)", realmFull, characterRealm)
	}

	validClassForRealm := false

	for _, class := range classesForRealm {
		if class == characterClass {
			validClassForRealm = true
		}
	}

	if !validClassForRealm {
		return fmt.Errorf("New character class %s (%d) not valid for realm %s (%d)",
			classFull, characterClass, realmFull, characterRealm)
	}

	classesForRace, ok := constants.CharacterRaceClasses[characterRace]
	if !ok {
		return fmt.Errorf("Unable to get classes for race %s (%d)", raceFull, characterRace)
	}

	validClassForRace := false

	for _, class := range classesForRace {
		if class == characterClass {
			validClassForRace = true
		}
	}

	if !validClassForRace {
		return fmt.Errorf("New character class %s (%d) not valid for race %s (%d)",
			classFull, characterClass, raceFull, characterRace)
	}

	if character.Gender > 0 {
		switch characterRace {
		case constants.CharacterRaceAlbionMinotaur:
		case constants.CharacterRaceMidgardMinotaur:
		case constants.CharacterRaceHiberniaMinotaur:
			return fmt.Errorf("New character gender %d is invalid for Minotaur race!", character.Gender)
		}
	}

	if character.Gender == 0 {
		switch characterClass {
		case constants.CharacterClassBainshee:
		case constants.CharacterClassValkyrie:
			return fmt.Errorf("New character gender %d is invalid for Bainshee/Valkyrie class!", character.Gender)
		}
	}

	pointsUsed, err := getPointsUsedOverBase(character)
	if err != nil {
		return err
	}

	if pointsUsed != 30 {
		return fmt.Errorf("New character must use all stat points!  Points used: %d", pointsUsed)
	}

	return nil
}

func getPointsUsedOverBase(character *models.Character) (int64, error) {
	max := func(a, b int64) int64 {
		if a > b {
			return a
		}

		return b
	}

	pointsForStat := func(base, current int64) int64 {
		current -= base

		result := current
		result += max(0, current-10)
		result += max(0, current-15)

		return result
	}

	sb := []uint32{
		character.Strength,
		character.Dexterity,
		character.Constitution,
		character.Quickness,
		character.Intelligence,
		character.Piety,
		character.Empathy,
		character.Charisma,
	}

	pointsUsed := int64(0)

	for i := 0; i < 8; i++ {
		baseValue := int64(constants.CharacterBaseRaceStats[constants.CharacterRace(character.Race)][i])
		currentValue := int64(sb[i])

		if currentValue < baseValue {
			return 0, fmt.Errorf("Starting value for %s less than base racial value! Base: %d, Given: %d", helpers.StatSlotToName(i), baseValue, currentValue)
		}

		if currentValue == baseValue {
			continue
		}

		pointsUsed += pointsForStat(baseValue, currentValue)
	}

	return pointsUsed, nil
}
