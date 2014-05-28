package interfaces

type Server interface {
	Config() Config
	Ruleset() Ruleset
	Database() Database
	World() World
	ScriptExecutor() ScriptExecutor

	ShortName() string
}
