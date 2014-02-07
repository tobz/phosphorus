package rulesets

import "fmt"
import "github.com/tobz/phosphorus/interfaces"

type rulesetCreator func(interfaces.Server) (interfaces.Ruleset, error)

type rulesetMap map[string]rulesetCreator

var rulesetCreators = make(rulesetMap)

func Register(rulesetName string, rulesetFunc rulesetCreator) {
	if _, ok := rulesetCreators[rulesetName]; ok {
		panic(fmt.Sprintf("ruleset called '%s'% already registered!", rulesetName))
	}

	rulesetCreators[rulesetName] = rulesetFunc
}

func GetRuleset(rulesetName string, s interfaces.Server) (interfaces.Ruleset, error) {
	f, ok := rulesetCreators[rulesetName]
	if !ok {
		return nil, fmt.Errorf("no ruleset named '%s' registered", rulesetName)
	}

	return f(s)
}
