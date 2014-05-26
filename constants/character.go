package constants

type CharacterClass uint32

const (
	CharacterClassAcolyte      CharacterClass = 16
	CharacterClassAlbionRogue  CharacterClass = 17
	CharacterClassDisciple     CharacterClass = 20
	CharacterClassElementalist CharacterClass = 15
	CharacterClassFighter      CharacterClass = 14
	CharacterClassForester     CharacterClass = 57
	CharacterClassGuardian     CharacterClass = 52
	CharacterClassMage         CharacterClass = 18
	CharacterClassMagician     CharacterClass = 51
	CharacterClassMidgardRogue CharacterClass = 38
	CharacterClassMystic       CharacterClass = 36
	CharacterClassNaturalist   CharacterClass = 53
	CharacterClassSeer         CharacterClass = 37
	CharacterClassStalker      CharacterClass = 54
	CharacterClassViking       CharacterClass = 35

	CharacterClassArmsman     CharacterClass = 2
	CharacterClassCabalist    CharacterClass = 13
	CharacterClassCleric      CharacterClass = 6
	CharacterClassFriar       CharacterClass = 10
	CharacterClassHeretic     CharacterClass = 33
	CharacterClassInfiltrator CharacterClass = 9
	CharacterClassMercenary   CharacterClass = 11
	CharacterClassMinstrel    CharacterClass = 4
	CharacterClassNecromancer CharacterClass = 12
	CharacterClassPaladin     CharacterClass = 1
	CharacterClassReaver      CharacterClass = 19
	CharacterClassScout       CharacterClass = 3
	CharacterClassSorcerer    CharacterClass = 8
	CharacterClassTheurgist   CharacterClass = 5
	CharacterClassWizard      CharacterClass = 7
	CharacterClassMaulerAlb   CharacterClass = 60

	CharacterClassBerserker    CharacterClass = 31
	CharacterClassBonedancer   CharacterClass = 30
	CharacterClassHealer       CharacterClass = 26
	CharacterClassHunter       CharacterClass = 25
	CharacterClassRunemaster   CharacterClass = 29
	CharacterClassSavage       CharacterClass = 32
	CharacterClassShadowblade  CharacterClass = 23
	CharacterClassShaman       CharacterClass = 28
	CharacterClassSkald        CharacterClass = 24
	CharacterClassSpiritmaster CharacterClass = 27
	CharacterClassThane        CharacterClass = 21
	CharacterClassValkyrie     CharacterClass = 34
	CharacterClassWarlock      CharacterClass = 59
	CharacterClassWarrior      CharacterClass = 22
	CharacterClassMaulerMid    CharacterClass = 61

	CharacterClassAnimist     CharacterClass = 55
	CharacterClassBainshee    CharacterClass = 39
	CharacterClassBard        CharacterClass = 48
	CharacterClassBlademaster CharacterClass = 43
	CharacterClassChampion    CharacterClass = 45
	CharacterClassDruid       CharacterClass = 47
	CharacterClassEldritch    CharacterClass = 40
	CharacterClassEnchanter   CharacterClass = 41
	CharacterClassHero        CharacterClass = 44
	CharacterClassMentalist   CharacterClass = 42
	CharacterClassNightshade  CharacterClass = 49
	CharacterClassRanger      CharacterClass = 50
	CharacterClassValewalker  CharacterClass = 56
	CharacterClassVampiir     CharacterClass = 58
	CharacterClassWarden      CharacterClass = 46
	CharacterClassMaulerHib   CharacterClass = 62
)

type CharacterRace uint32

const (
	CharacterRaceBriton           CharacterRace = 1
	CharacterRaceAvalonian        CharacterRace = 2
	CharacterRaceHighlander       CharacterRace = 3
	CharacterRaceSaracen          CharacterRace = 4
	CharacterRaceNorseman         CharacterRace = 5
	CharacterRaceTroll            CharacterRace = 6
	CharacterRaceDwarf            CharacterRace = 7
	CharacterRaceKobold           CharacterRace = 8
	CharacterRaceCelt             CharacterRace = 9
	CharacterRaceFirbolg          CharacterRace = 10
	CharacterRaceElf              CharacterRace = 11
	CharacterRaceLurikeen         CharacterRace = 12
	CharacterRaceInconnu          CharacterRace = 13
	CharacterRaceValkyn           CharacterRace = 14
	CharacterRaceSylvan           CharacterRace = 15
	CharacterRaceHalfOgre         CharacterRace = 16
	CharacterRaceFrostalf         CharacterRace = 17
	CharacterRaceShar             CharacterRace = 18
	CharacterRaceAlbionMinotaur   CharacterRace = 19
	CharacterRaceMidgardMinotaur  CharacterRace = 20
	CharacterRaceHiberniaMinotaur CharacterRace = 21
)

var CharacterStartingClasses map[ClientRealm][]CharacterClass = map[ClientRealm][]CharacterClass{
	ClientRealmAlbion: []CharacterClass{
		CharacterClassPaladin,
		CharacterClassArmsman,
		CharacterClassScout,
		CharacterClassMinstrel,
		CharacterClassTheurgist,
		CharacterClassCleric,
		CharacterClassWizard,
		CharacterClassSorcerer,
		CharacterClassInfiltrator,
		CharacterClassFriar,
		CharacterClassMercenary,
		CharacterClassNecromancer,
		CharacterClassCabalist,
		CharacterClassFighter,
		CharacterClassElementalist,
		CharacterClassAcolyte,
		CharacterClassAlbionRogue,
		CharacterClassMage,
		CharacterClassReaver,
		CharacterClassDisciple,
		CharacterClassHeretic,
		CharacterClassMaulerAlb,
	},
	ClientRealmMidgard: []CharacterClass{
		CharacterClassThane,
		CharacterClassWarrior,
		CharacterClassShadowblade,
		CharacterClassSkald,
		CharacterClassHunter,
		CharacterClassHealer,
		CharacterClassSpiritmaster,
		CharacterClassShaman,
		CharacterClassRunemaster,
		CharacterClassBonedancer,
		CharacterClassBerserker,
		CharacterClassSavage,
		CharacterClassValkyrie,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
		CharacterClassMidgardRogue,
		CharacterClassWarlock,
		CharacterClassMaulerMid,
	},
	ClientRealmHibernia: []CharacterClass{
		CharacterClassBainshee,
		CharacterClassEldritch,
		CharacterClassEnchanter,
		CharacterClassMentalist,
		CharacterClassBlademaster,
		CharacterClassHero,
		CharacterClassChampion,
		CharacterClassWarden,
		CharacterClassDruid,
		CharacterClassBard,
		CharacterClassNightshade,
		CharacterClassRanger,
		CharacterClassMagician,
		CharacterClassGuardian,
		CharacterClassNaturalist,
		CharacterClassStalker,
		CharacterClassAnimist,
		CharacterClassValewalker,
		CharacterClassForester,
		CharacterClassVampiir,
		CharacterClassMaulerHib,
	},
}

var CharacterRaceClasses map[CharacterRace][]CharacterClass = map[CharacterRace][]CharacterClass{
	CharacterRaceBriton: []CharacterClass{
		CharacterClassArmsman,
		CharacterClassReaver,
		CharacterClassMercenary,
		CharacterClassPaladin,
		CharacterClassCleric,
		CharacterClassHeretic,
		CharacterClassFriar,
		CharacterClassSorcerer,
		CharacterClassCabalist,
		CharacterClassTheurgist,
		CharacterClassNecromancer,
		CharacterClassMaulerAlb,
		CharacterClassWizard,
		CharacterClassMinstrel,
		CharacterClassInfiltrator,
		CharacterClassScout,
		CharacterClassFighter,
		CharacterClassAcolyte,
		CharacterClassMage,
		CharacterClassElementalist,
		CharacterClassAlbionRogue,
		CharacterClassDisciple,
	},
	CharacterRaceAvalonian: []CharacterClass{
		CharacterClassPaladin,
		CharacterClassCleric,
		CharacterClassWizard,
		CharacterClassTheurgist,
		CharacterClassArmsman,
		CharacterClassMercenary,
		CharacterClassSorcerer,
		CharacterClassCabalist,
		CharacterClassHeretic,
		CharacterClassFriar,
		CharacterClassFighter,
		CharacterClassAcolyte,
		CharacterClassMage,
		CharacterClassElementalist,
	},
	CharacterRaceHighlander: []CharacterClass{
		CharacterClassArmsman,
		CharacterClassMercenary,
		CharacterClassPaladin,
		CharacterClassCleric,
		CharacterClassMinstrel,
		CharacterClassScout,
		CharacterClassFriar,
		CharacterClassFighter,
		CharacterClassAcolyte,
		CharacterClassAlbionRogue,
	},
	CharacterRaceSaracen: []CharacterClass{
		CharacterClassSorcerer,
		CharacterClassCabalist,
		CharacterClassPaladin,
		CharacterClassReaver,
		CharacterClassMercenary,
		CharacterClassArmsman,
		CharacterClassInfiltrator,
		CharacterClassMinstrel,
		CharacterClassScout,
		CharacterClassNecromancer,
		CharacterClassFighter,
		CharacterClassMage,
		CharacterClassAlbionRogue,
		CharacterClassDisciple,
	},
	CharacterRaceNorseman: []CharacterClass{
		CharacterClassHealer,
		CharacterClassWarrior,
		CharacterClassBerserker,
		CharacterClassThane,
		CharacterClassWarlock,
		CharacterClassSkald,
		CharacterClassValkyrie,
		CharacterClassSpiritmaster,
		CharacterClassRunemaster,
		CharacterClassSavage,
		CharacterClassMaulerMid,
		CharacterClassShadowblade,
		CharacterClassHunter,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
		CharacterClassMidgardRogue,
	},
	CharacterRaceTroll: []CharacterClass{
		CharacterClassBerserker,
		CharacterClassWarrior,
		CharacterClassSavage,
		CharacterClassThane,
		CharacterClassSkald,
		CharacterClassBonedancer,
		CharacterClassShaman,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
	},
	CharacterRaceDwarf: []CharacterClass{
		CharacterClassHealer,
		CharacterClassThane,
		CharacterClassBerserker,
		CharacterClassWarrior,
		CharacterClassSavage,
		CharacterClassSkald,
		CharacterClassValkyrie,
		CharacterClassRunemaster,
		CharacterClassHunter,
		CharacterClassShaman,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
		CharacterClassMidgardRogue,
	},
	CharacterRaceKobold: []CharacterClass{
		CharacterClassShaman,
		CharacterClassWarrior,
		CharacterClassSkald,
		CharacterClassSavage,
		CharacterClassRunemaster,
		CharacterClassSpiritmaster,
		CharacterClassBonedancer,
		CharacterClassWarlock,
		CharacterClassHunter,
		CharacterClassShadowblade,
		CharacterClassMaulerMid,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
		CharacterClassMidgardRogue,
	},
	CharacterRaceCelt: []CharacterClass{
		CharacterClassBard,
		CharacterClassDruid,
		CharacterClassWarden,
		CharacterClassBlademaster,
		CharacterClassHero,
		CharacterClassVampiir,
		CharacterClassChampion,
		CharacterClassMaulerHib,
		CharacterClassMentalist,
		CharacterClassBainshee,
		CharacterClassRanger,
		CharacterClassAnimist,
		CharacterClassValewalker,
		CharacterClassNightshade,
		CharacterClassGuardian,
		CharacterClassStalker,
		CharacterClassNaturalist,
		CharacterClassMagician,
		CharacterClassForester,
	},
	CharacterRaceFirbolg: []CharacterClass{
		CharacterClassBard,
		CharacterClassDruid,
		CharacterClassWarden,
		CharacterClassHero,
		CharacterClassBlademaster,
		CharacterClassAnimist,
		CharacterClassValewalker,
		CharacterClassGuardian,
		CharacterClassNaturalist,
		CharacterClassForester,
	},
	CharacterRaceElf: []CharacterClass{
		CharacterClassBlademaster,
		CharacterClassChampion,
		CharacterClassRanger,
		CharacterClassNightshade,
		CharacterClassBainshee,
		CharacterClassEnchanter,
		CharacterClassEldritch,
		CharacterClassMentalist,
		CharacterClassGuardian,
		CharacterClassStalker,
		CharacterClassMagician,
	},
	CharacterRaceLurikeen: []CharacterClass{
		CharacterClassHero,
		CharacterClassChampion,
		CharacterClassVampiir,
		CharacterClassEldritch,
		CharacterClassEnchanter,
		CharacterClassMentalist,
		CharacterClassBainshee,
		CharacterClassNightshade,
		CharacterClassRanger,
		CharacterClassMaulerHib,
		CharacterClassGuardian,
		CharacterClassStalker,
		CharacterClassMagician,
	},
	CharacterRaceInconnu: []CharacterClass{
		CharacterClassReaver,
		CharacterClassSorcerer,
		CharacterClassCabalist,
		CharacterClassHeretic,
		CharacterClassNecromancer,
		CharacterClassArmsman,
		CharacterClassMercenary,
		CharacterClassInfiltrator,
		CharacterClassScout,
		CharacterClassMaulerAlb,
		CharacterClassFighter,
		CharacterClassAcolyte,
		CharacterClassMage,
		CharacterClassAlbionRogue,
		CharacterClassDisciple,
	},
	CharacterRaceValkyn: []CharacterClass{
		CharacterClassSavage,
		CharacterClassBerserker,
		CharacterClassBonedancer,
		CharacterClassWarrior,
		CharacterClassShadowblade,
		CharacterClassHunter,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassMidgardRogue,
	},
	CharacterRaceSylvan: []CharacterClass{
		CharacterClassAnimist,
		CharacterClassDruid,
		CharacterClassValewalker,
		CharacterClassHero,
		CharacterClassWarden,
		CharacterClassGuardian,
		CharacterClassNaturalist,
		CharacterClassForester,
	},
	CharacterRaceHalfOgre: []CharacterClass{
		CharacterClassWizard,
		CharacterClassTheurgist,
		CharacterClassCabalist,
		CharacterClassSorcerer,
		CharacterClassMercenary,
		CharacterClassArmsman,
		CharacterClassFighter,
		CharacterClassMage,
		CharacterClassElementalist,
	},
	CharacterRaceFrostalf: []CharacterClass{
		CharacterClassHealer,
		CharacterClassShaman,
		CharacterClassThane,
		CharacterClassSpiritmaster,
		CharacterClassRunemaster,
		CharacterClassWarlock,
		CharacterClassValkyrie,
		CharacterClassHunter,
		CharacterClassShadowblade,
		CharacterClassViking,
		CharacterClassMystic,
		CharacterClassSeer,
		CharacterClassMidgardRogue,
	},
	CharacterRaceShar: []CharacterClass{
		CharacterClassChampion,
		CharacterClassHero,
		CharacterClassBlademaster,
		CharacterClassVampiir,
		CharacterClassRanger,
		CharacterClassMentalist,
		CharacterClassGuardian,
		CharacterClassStalker,
		CharacterClassMagician,
	},
	CharacterRaceAlbionMinotaur: []CharacterClass{
		CharacterClassHeretic,
		CharacterClassMaulerAlb,
		CharacterClassArmsman,
		CharacterClassMercenary,
		CharacterClassFighter,
		CharacterClassAcolyte,
	},
	CharacterRaceMidgardMinotaur: []CharacterClass{
		CharacterClassBerserker,
		CharacterClassMaulerMid,
		CharacterClassThane,
		CharacterClassViking,
		CharacterClassWarrior,
	},
	CharacterRaceHiberniaMinotaur: []CharacterClass{
		CharacterClassHero,
		CharacterClassBlademaster,
		CharacterClassMaulerHib,
		CharacterClassWarden,
		CharacterClassGuardian,
		CharacterClassNaturalist,
	},
}

// Stats are ordered as follows: strength, dexterity, constitution, quickness,
// intelligence, piety, empathy, and charisma.
var CharacterBaseRaceStats map[CharacterRace][]uint8 = map[CharacterRace][]uint8{
	CharacterRaceBriton:           []uint8{60, 60, 60, 60, 60, 60, 60, 60},
	CharacterRaceAvalonian:        []uint8{45, 60, 45, 70, 80, 60, 60, 60},
	CharacterRaceHighlander:       []uint8{70, 50, 70, 50, 60, 60, 60, 60},
	CharacterRaceSaracen:          []uint8{50, 80, 50, 60, 60, 60, 60, 60},
	CharacterRaceNorseman:         []uint8{70, 50, 70, 50, 60, 60, 60, 60},
	CharacterRaceTroll:            []uint8{100, 35, 70, 35, 60, 60, 60, 60},
	CharacterRaceDwarf:            []uint8{60, 50, 80, 50, 60, 60, 60, 60},
	CharacterRaceKobold:           []uint8{50, 70, 50, 70, 60, 60, 60, 60},
	CharacterRaceCelt:             []uint8{60, 60, 60, 60, 60, 60, 60, 60},
	CharacterRaceFirbolg:          []uint8{90, 40, 60, 40, 60, 60, 70, 60},
	CharacterRaceElf:              []uint8{40, 75, 40, 75, 70, 60, 60, 60},
	CharacterRaceLurikeen:         []uint8{40, 80, 40, 80, 60, 60, 60, 60},
	CharacterRaceInconnu:          []uint8{50, 70, 60, 50, 70, 60, 60, 60},
	CharacterRaceValkyn:           []uint8{55, 65, 45, 75, 60, 60, 60, 60},
	CharacterRaceSylvan:           []uint8{70, 55, 60, 45, 70, 60, 60, 60},
	CharacterRaceHalfOgre:         []uint8{90, 40, 70, 40, 60, 60, 60, 60},
	CharacterRaceFrostalf:         []uint8{55, 55, 55, 60, 60, 75, 60, 60},
	CharacterRaceShar:             []uint8{60, 50, 80, 50, 60, 60, 60, 60},
	CharacterRaceAlbionMinotaur:   []uint8{80, 50, 70, 40, 60, 60, 60, 60},
	CharacterRaceMidgardMinotaur:  []uint8{80, 50, 70, 40, 60, 60, 60, 60},
	CharacterRaceHiberniaMinotaur: []uint8{80, 50, 70, 40, 60, 60, 60, 60},
}

const (
	BaseCharacterSpeed = 191
)
