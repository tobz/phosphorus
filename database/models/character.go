package models

import "time"

type Character struct {
	CharacterID uint64 `db:"character_id"`
	AccountID   uint64 `db:"account_id"`
	AccountSlot uint32 `db:"account_slot"`
	Realm       uint8  `db:"realm"`

	Created    time.Time `db:"created_dt"`
	LastPlayed time.Time `db:"last_played_dt"`

	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	Race              uint8  `db:"race"`
	Class             uint8  `db:"class"`
	Gender            uint8  `db:"gender"`
	EyeSize           uint8  `db:"eye_size"`
	LipSize           uint8  `db:"lip_size"`
	EyeColor          uint8  `db:"eye_color"`
	HairColor         uint8  `db:"hair_color"`
	FaceType          uint8  `db:"face_type"`
	HairStyle         uint8  `db:"hair_style"`
	MoodType          uint8  `db:"mood_type"`
	BaseModel         uint16 `db:"base_model"`
	CurrentModel      uint16 `db:"current_model"`
	CustomizationStep uint8  `db:"customization_step"`

	Constitution uint32 `db:"constitution"`
	Dexterity    uint32 `db:"dexterity"`
	Strength     uint32 `db:"strength"`
	Quickness    uint32 `db:"quickness"`
	Intelligence uint32 `db:"intelligence"`
	Piety        uint32 `db:"piety"`
	Empathy      uint32 `db:"empathy"`
	Charisma     uint32 `db:"charisma"`

	Level      uint8  `db:"level"`
	Experience uint64 `db:"experience"`

	Endurance    uint32 `db:"endurance"`
	MaxEndurance uint32 `db:"max_endurance"`

	MaxSpeed uint32 `db:"max_speed"`

	PositionX uint32 `db:"position_x"`
	PositionY uint32 `db:"position_y"`
	PositionZ uint32 `db:"position_z"`
	Region    uint8  `db:"region"`
	Heading   uint32 `db:"heading"`

	BindPositionX uint32 `db:"bind_position_x"`
	BindPositionY uint32 `db:"bind_position_y"`
	BindPositionZ uint32 `db:"bind_position_z"`
	BindRegion    uint32 `db:"bind_region"`
	BindHeading   uint32 `db:"bind_heading"`

	GuildID uint32 `db:"guild_id"`
}
