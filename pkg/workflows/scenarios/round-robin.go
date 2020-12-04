package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"time"
)

type RoundRobin struct{}

func (r RoundRobin) Generate(presetsList presets.List, params ScenarioParams) (Scenario, error) {
	if len(presetsList.ContainerPresets)+len(presetsList.PodPresets) == 0 {
		return Scenario{}, ZeroActions
	}

	if params.Stages <= 0 {
		return Scenario{}, NonPositiveStagesError
	}

	if params.Stages > 100 {
		return Scenario{}, TooManyStagesError
	}

	stages := make([]Stage, 0, params.Stages)

	stagesLeft := params.Stages
	for stagesLeft > 0 {
		for _, preset := range presetsList.ContainerPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			newAction := PlannedAction{Type: preset.Type(), Engine: preset.Instantiate("app=server", "server")}
			stage := Stage{Actions: []PlannedAction{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}

		for _, preset := range presetsList.PodPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			newAction := PlannedAction{Type: preset.Type(), Engine: preset.Instantiate("app=server")}
			stage := Stage{Actions: []PlannedAction{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}
	}

	return Scenario{stages}, nil
}

func NewRoundRobinGenerator() RoundRobin {
	return RoundRobin{}
}
