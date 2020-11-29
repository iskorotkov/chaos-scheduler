package scenarios

import (
	"errors"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroActions            = errors.New("can't create scenario out of 0 actions")
)

type RoundRobin byte

func (r RoundRobin) Generate(actions []ActionTemplate, config Config) (Scenario, error) {
	if len(actions) == 0 {
		return nil, ZeroActions
	}

	if config.Stages <= 0 {
		return nil, NonPositiveStagesError
	}

	if config.Stages > 100 {
		return nil, TooManyStagesError
	}

	scenario := make([]Stage, 0, config.Stages)

	for i := 0; i < config.Stages; i++ {
		a := actions[i%len(actions)]

		stage := []PlannedAction{{Name: a.Name, Template: a.Template}}
		scenario = append(scenario, stage)
	}

	return scenario, nil
}

func NewRoundRobinGenerator() Generator {
	return RoundRobin(0)
}
