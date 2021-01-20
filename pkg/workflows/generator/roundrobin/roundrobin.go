package roundrobin

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
)

type RoundRobin struct {
	presets []experiments.Preset
	seeker  targets.Seeker
	logger  *zap.SugaredLogger
}

func (r RoundRobin) Generate(params generator.Params) (generator.Scenario, error) {
	if len(r.presets) == 0 {
		return generator.Scenario{}, generator.ZeroFailures
	}

	if params.Stages <= 0 {
		return generator.Scenario{}, generator.NonPositiveStagesError
	}

	if params.Stages > 100 {
		return generator.Scenario{}, generator.TooManyStagesError
	}

	src := rand.NewSource(params.Seed)
	rnd := rand.New(src)

	targetsList, err := r.seeker.Targets()
	if err != nil {
		r.logger.Error(err.Error())
		return generator.Scenario{}, generator.TargetsError
	}

	stages := make([]generator.Stage, 0, params.Stages)

	stagesLeft := params.Stages
	for stagesLeft > 0 {
		for _, preset := range r.presets {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			target := selectTarget(targetsList, rnd)
			engine := preset.Engine(target, params.StageDuration)
			newAction := generator.Action{
				Info:   preset.Info(),
				Target: target,
				Engine: engine,
			}

			stage := generator.Stage{Actions: []generator.Action{newAction}, Duration: params.StageDuration}
			stages = append(stages, stage)
		}
	}

	return generator.Scenario{Stages: stages}, nil
}

func NewRoundRobin(presets []experiments.Preset, seeker targets.KubernetesSeeker, logger *zap.SugaredLogger) RoundRobin {
	return RoundRobin{presets: presets, seeker: seeker, logger: logger}
}

func selectTarget(ts []targets.Target, rnd *rand.Rand) targets.Target {
	index := rnd.Intn(len(ts))
	return ts[index]
}
