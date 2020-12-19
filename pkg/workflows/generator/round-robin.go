package generator

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
)

type RoundRobin struct {
	presets []experiments.Preset
	seeker  targets.Seeker
	logger  *zap.SugaredLogger
}

func (r RoundRobin) Generate(params Params) (Scenario, error) {
	if len(r.presets) == 0 {
		return Scenario{}, ZeroFailures
	}

	if params.Stages <= 0 {
		return Scenario{}, NonPositiveStagesError
	}

	if params.Stages > 100 {
		return Scenario{}, TooManyStagesError
	}

	src := rand.NewSource(params.Seed)
	rnd := rand.New(src)

	targetsList, err := r.seeker.Targets()
	if err != nil {
		r.logger.Error(err.Error())
		return Scenario{}, TargetsError
	}

	stages := make([]Stage, 0, params.Stages)

	stagesLeft := params.Stages
	for stagesLeft > 0 {
		for _, preset := range r.presets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			target := selectTarget(targetsList, rnd)
			engine := preset.Engine(target, params.StageDuration)
			newAction := Action{
				Info:   preset.Info(),
				Target: target,
				Engine: engine,
			}

			stage := Stage{Actions: []Action{newAction}, Duration: params.StageDuration}
			stages = append(stages, stage)
		}
	}

	return Scenario{stages}, nil
}

func NewRoundRobin(presets []experiments.Preset, seeker targets.KubernetesSeeker, logger *zap.SugaredLogger) RoundRobin {
	return RoundRobin{presets: presets, seeker: seeker, logger: logger}
}

func selectTarget(ts []targets.Target, rnd *rand.Rand) targets.Target {
	index := rnd.Intn(len(ts))
	return ts[index]
}
