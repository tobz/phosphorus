package rulesets

import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/interfaces"

func init() {
	Register("rvr", NewRvRRuleset)
}

type RvRRuleset struct {
	server interfaces.Server
}

func NewRvRRuleset(s interfaces.Server) (interfaces.Ruleset, error) {
	return &RvRRuleset{s}, nil
}

func (rs *RvRRuleset) ColorHandling() constants.ColorHandling {
	return constants.ColorHandlingRvR
}
