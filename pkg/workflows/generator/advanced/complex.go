package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func addComplexFailures(a *Generator, targetsList []targets.Target, r *rand.Rand, params phaseParams) []generator.Stage {
	stages := make([]generator.Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageFailures := make([]Failure, len(a.failures))
		copy(stageFailures, a.failures)

		stageTargets := make([]targets.Target, len(targetsList))
		copy(stageTargets, targetsList)

		actions := make([]generator.Action, 0)

		points := a.budget.MaxPoints
		retries := a.retries
		for len(actions) < a.budget.MaxFailures {
			failure := popRandomFailure(stageFailures, r)
			target := popRandomTarget(stageTargets, r)
			cost := calculateCost(a.modifiers, failure)

			if cost <= points {
				points -= cost

				actions = append(actions, generator.Action{
					Info:   failure.Preset.Info(),
					Target: target,
					Engine: failure.Preset.Engine(target, params.StageDuration),
				})
			} else {
				if retries <= 0 {
					break
				}

				retries--
			}
		}
	}

	return stages
}
