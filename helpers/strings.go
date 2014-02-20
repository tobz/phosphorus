package helpers

import "strings"
import "github.com/tobz/phosphorus/constants"

func GetClassName(class uint8, capitalize bool) string {
	classConst := constants.CharacterClass(class)
	className := ""

	switch classConst {
	case constants.CharacterClassAcolyte:
		className = "acolyte"
	case constants.CharacterClassAlbionRogue:
		className = "rogue"
	case constants.CharacterClassDisciple:
		className = "disciple"
	case constants.CharacterClassElementalist:
		className = "elementalist"
	case constants.CharacterClassFighter:
		className = "fighter"
	case constants.CharacterClassForester:
		className = "forester"
	case constants.CharacterClassGuardian:
		className = "guardian"
	case constants.CharacterClassMage:
		className = "mage"
	case constants.CharacterClassMagician:
		className = "magician"
	case constants.CharacterClassMidgardRogue:
		className = "rogue"
	case constants.CharacterClassMystic:
		className = "mystic"
	case constants.CharacterClassNaturalist:
		className = "naturalist"
	case constants.CharacterClassSeer:
		className = "seer"
	case constants.CharacterClassStalker:
		className = "stalker"
	case constants.CharacterClassViking:
		className = "viking"

	// Albion classes.
	case constants.CharacterClassArmsman:
		className = "armsman"
	case constants.CharacterClassCabalist:
		className = "cabalist"
	case constants.CharacterClassCleric:
		className = "cleric"
	case constants.CharacterClassFriar:
		className = "friar"
	case constants.CharacterClassHeretic:
		className = "heretic"
	case constants.CharacterClassInfiltrator:
		className = "infiltrator"
	case constants.CharacterClassMercenary:
		className = "mercenary"
	case constants.CharacterClassMinstrel:
		className = "minstrel"
	case constants.CharacterClassNecromancer:
		className = "necromancer"
	case constants.CharacterClassPaladin:
		className = "paladin"
	case constants.CharacterClassReaver:
		className = "reaver"
	case constants.CharacterClassScout:
		className = "scout"
	case constants.CharacterClassSorcerer:
		className = "sorcerer"
	case constants.CharacterClassTheurgist:
		className = "theurgist"
	case constants.CharacterClassWizard:
		className = "wizard"
	case constants.CharacterClassMaulerAlb:
		className = "mauler"

	// Midgard classes.
	case constants.CharacterClassBerserker:
		className = "berserker"
	case constants.CharacterClassBonedancer:
		className = "bonedancer"
	case constants.CharacterClassHealer:
		className = "healer"
	case constants.CharacterClassHunter:
		className = "hunter"
	case constants.CharacterClassRunemaster:
		className = "runemaster"
	case constants.CharacterClassSavage:
		className = "savage"
	case constants.CharacterClassShadowblade:
		className = "shadowblade"
	case constants.CharacterClassShaman:
		className = "shaman"
	case constants.CharacterClassSkald:
		className = "skald"
	case constants.CharacterClassSpiritmaster:
		className = "spiritmaster"
	case constants.CharacterClassThane:
		className = "thane"
	case constants.CharacterClassValkyrie:
		className = "valkyrie"
	case constants.CharacterClassWarlock:
		className = "warlock"
	case constants.CharacterClassWarrior:
		className = "warrior"
	case constants.CharacterClassMaulerMid:
		className = "mauler"

	// Hibernia classes.
	case constants.CharacterClassAnimist:
		className = "animist"
	case constants.CharacterClassBainshee:
		className = "bainshee"
	case constants.CharacterClassBard:
		className = "bard"
	case constants.CharacterClassBlademaster:
		className = "blademaster"
	case constants.CharacterClassChampion:
		className = "champion"
	case constants.CharacterClassDruid:
		className = "druid"
	case constants.CharacterClassEldritch:
		className = "eldritch"
	case constants.CharacterClassEnchanter:
		className = "enchanter"
	case constants.CharacterClassHero:
		className = "hero"
	case constants.CharacterClassMentalist:
		className = "mentalist"
	case constants.CharacterClassNightshade:
		className = "nightshade"
	case constants.CharacterClassRanger:
		className = "ranger"
	case constants.CharacterClassValewalker:
		className = "valewalker"
	case constants.CharacterClassVampiir:
		className = "vampiir"
	case constants.CharacterClassWarden:
		className = "warden"
	case constants.CharacterClassMaulerHib:
		className = "mauler"
	}

	if capitalize {
		return strings.Title(className)
	}

	return className
}

func GetRaceName(race uint8, capitalize bool) string {
	raceConst := constants.CharacterRace(race)
	raceName := ""

	switch raceConst {
	case constants.CharacterRaceTroll:
		raceName = "troll"
	case constants.CharacterRaceDwarf:
		raceName = "dwarf"
	case constants.CharacterRaceKobold:
		raceName = "kobold"
	case constants.CharacterRaceCelt:
		raceName = "celt"
	case constants.CharacterRaceFirbolg:
		raceName = "firbolg"
	case constants.CharacterRaceElf:
		raceName = "elf"
	case constants.CharacterRaceLurikeen:
		raceName = "lurikeen"
	case constants.CharacterRaceInconnu:
		raceName = "inconnu"
	case constants.CharacterRaceValkyn:
		raceName = "valkyn"
	case constants.CharacterRaceSylvan:
		raceName = "sylvan"
	case constants.CharacterRaceHalfOgre:
		raceName = "half ogre"
	case constants.CharacterRaceFrostalf:
		raceName = "frostalf"
	case constants.CharacterRaceShar:
		raceName = "shar"
	case constants.CharacterRaceAlbionMinotaur:
	case constants.CharacterRaceMidgardMinotaur:
	case constants.CharacterRaceHiberniaMinotaur:
		raceName = "minotaur"
	}

	if capitalize {
		return strings.Title(raceName)
	}

	return raceName
}
