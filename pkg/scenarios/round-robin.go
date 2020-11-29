package scenarios

import (
	"time"
)

type RoundRobin byte

func (r RoundRobin) Generate(actions []Template, config Config) (Scenario, error) {
	if len(actions) == 0 {
		return nil, ZeroActions
	}

	if config.Stages <= 0 {
		return nil, NonPositiveStagesError
	}

	if config.Stages > 100 {
		return nil, TooManyStagesError
	}

	stages := make([]Stage, 0, config.Stages)

	for i := 0; i < config.Stages; i++ {
		a := actions[i%len(actions)]

		stage := stage{action{name: a.Name(), template: a.Template(), duration: time.Minute}}
		stages = append(stages, stage)
	}

	return scenario(stages), nil
}

func NewRoundRobinGenerator() Generator {
	return RoundRobin(0)
}
