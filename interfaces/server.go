package interfaces

type Server interface {
	Config() Config
	Ruleset() Ruleset
	Database() Database
	World() World

	ShortName() string
}
