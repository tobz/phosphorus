package helpers

import "github.com/tobz/phosphorus/constants"

func StatSlotToName(statSlot int) string {
	switch statSlot {
	case 0:
		return "strength"
	case 1:
		return "dexterity"
	case 2:
		return "constitution"
	case 3:
		return "quickness"
	case 4:
		return "intelligence"
	case 5:
		return "piety"
	case 6:
		return "empathy"
	case 7:
		return "charisma"
	default:
		return "unknown"
	}
}

func RealmToString(realm constants.ClientRealm) string {
	switch realm {
	case constants.ClientRealmAlbion:
		return "Albion"
	case constants.ClientRealmMidgard:
		return "Midgard"
	case constants.ClientRealmHibernia:
		return "Hibernia"
	default:
		return "unknown"
	}
}

func ClassToString(class constants.CharacterClass) string {
	return "unknown"
}

func RaceToString(race constants.CharacterRace) string {
	return "unknown"
}
