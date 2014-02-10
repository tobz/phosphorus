package interfaces

type Server interface {
	Config() Config
	Ruleset() Ruleset
    Database() Database

	ShortName() string
}
