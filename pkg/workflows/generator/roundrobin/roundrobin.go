package roundrobin

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
)

type RoundRobin struct {
	failures []failures.Failure
	seeker   targets.Seeker
	logger   *zap.SugaredLogger
}

func (r RoundRobin) Generate(params generator.Params) (generator.Scenario, error) {
	if len(r.failures) == 0 {
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
		for _, failure := range r.failures {
			if stagesLeft == 0 {
				break
			}

			stagesLeft--

			target := selectTarget(targetsList, rnd)
			engine := failure.Template.Instantiate(target, params.StageDuration)
			newAction := generator.Action{
				Name:     failure.Name(),
				Severity: failure.Severity,
				Scale:    failure.Scale,
				Target:   target,
				Engine:   engine,
			}

			stage := generator.Stage{Actions: []generator.Action{newAction}, Duration: params.StageDuration}
			stages = append(stages, stage)
		}
	}

	return generator.Scenario{Stages: stages}, nil
}

func NewRoundRobin(failures []failures.Failure, seeker targets.Seeker, logger *zap.SugaredLogger) RoundRobin {
	return RoundRobin{failures: failures, seeker: seeker, logger: logger}
}

func selectTarget(ts []targets.Target, rnd *rand.Rand) targets.Target {
	index := rnd.Intn(len(ts))
	return ts[index]
}
