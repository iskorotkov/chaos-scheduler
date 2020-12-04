package generators

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"math/rand"
	"time"
)

var (
	TargetsError = errors.New("couldn't get list of targets")
)

type RoundRobin struct {
	Presets presets.List
	Seeker  targets.Seeker
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

	src := rand.NewSource(params.Seed)
	rnd := rand.New(src)

	targetsList, err := r.Seeker.Targets()
	if err != nil {
		logger.Error(err)
		return Scenario{}, TargetsError
	}

	stages := make([]Stage, 0, params.Stages)

	stagesLeft := params.Stages
	for stagesLeft > 0 {
		for _, preset := range r.Presets.ContainerPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			target := selectTarget(targetsList, rnd)
			engine := preset.Instantiate(target.Selector(), target.MainContainer())
			newAction := Action{Type: preset.Type(), Engine: engine}

			stage := Stage{Actions: []Action{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}

		for _, preset := range r.Presets.PodPresets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			target := selectTarget(targetsList, rnd)
			engine := preset.Instantiate(target.Selector())
			newAction := Action{Type: preset.Type(), Engine: engine}

			stage := Stage{Actions: []Action{newAction}, Duration: time.Minute}
			stages = append(stages, stage)
		}
	}

	return Scenario{stages}, nil
}

func NewRoundRobinGenerator(presetsList presets.List, seeker targets.Seeker) RoundRobin {
	return RoundRobin{Presets: presetsList, Seeker: seeker}
}

func selectTarget(ts []targets.Target, rnd *rand.Rand) targets.Target {
	index := rnd.Intn(len(ts))
	return ts[index]
}
