package interfaces

type Server interface {
    Config() Config
    Ruleset() Ruleset

    ShortName() string
}
