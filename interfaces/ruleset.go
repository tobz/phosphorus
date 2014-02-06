package interfaces

import "github.com/tobz/phosphorus/constants"

type Ruleset interface {
    ColorHandling() constants.ColorHandling
}
