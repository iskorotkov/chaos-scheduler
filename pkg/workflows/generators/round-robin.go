package generators

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"time"
)

type RoundRobin struct {
	Presets presets.List
}

func (r RoundRobin) Generate(params Params) (Scenario, error) {
	if len(r.Presets.ContainerPresets)+len(r.Presets.PodPresets) == 0 {
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
		for _, preset := range r.Presets.ContainerPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			newAction := Action{Type: preset.Type(), Engine: preset.Instantiate("app=server", "server")}
			stage := Stage{Actions: []Action{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}

		for _, preset := range r.Presets.PodPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			newAction := Action{Type: preset.Type(), Engine: preset.Instantiate("app=server")}
			stage := Stage{Actions: []Action{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}
	}

	return Scenario{stages}, nil
}

func NewRoundRobinGenerator(presetsList presets.List) RoundRobin {
	return RoundRobin{Presets: presetsList}
}
