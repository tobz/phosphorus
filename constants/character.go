package constants

type CharacterClass uint32

const (
	// Base classes.
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

	// Albion classes.
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

	// Midgard classes.
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

	// Hibernia classes.
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
