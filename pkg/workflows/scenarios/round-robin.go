package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"
	"time"
)

type RoundRobin struct{}

func (r RoundRobin) Generate(factories []engines.Factory, params ScenarioParams) (Scenario, error) {
	if len(factories) == 0 {
		return Scenario{}, ZeroActions
	}

	if params.Stages <= 0 {
		return Scenario{}, NonPositiveStagesError
	}

	if params.Stages > 100 {
		return Scenario{}, TooManyStagesError
	}

	stages := make([]Stage, 0, params.Stages)

	for i := 0; i < params.Stages; i++ {
		factory := factories[i%len(factories)]
		engine := factory.Create("server")

		newAction := PlannedAction{Type: factory.Type(), Engine: engine}
		stage := Stage{Actions: []PlannedAction{newAction}, Duration: time.Minute}
		stages = append(stages, stage)
	}

	return Scenario{stages}, nil
}

func NewRoundRobinGenerator() RoundRobin {
	return RoundRobin{}
}
