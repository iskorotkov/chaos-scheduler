package scenarios

import (
	"time"
)

type RoundRobin byte

func (r RoundRobin) Generate(actions []TemplatedAction, config Config) (Scenario, error) {
	if len(actions) == 0 {
		return Scenario{}, ZeroActions
	}

	if config.Stages <= 0 {
		return Scenario{}, NonPositiveStagesError
	}

	if config.Stages > 100 {
		return Scenario{}, TooManyStagesError
	}

	stages := make([]Stage, 0, config.Stages)

	for i := 0; i < config.Stages; i++ {
		actionTemplate := actions[i%len(actions)]

		newAction := PlannedAction{Name: actionTemplate.Name, Template: actionTemplate.Template}
		stage := Stage{Actions: []PlannedAction{newAction}, Duration: time.Minute}
		stages = append(stages, stage)
	}

	return Scenario{stages}, nil
}

func NewRoundRobinGenerator() Generator {
	return RoundRobin(0)
}
