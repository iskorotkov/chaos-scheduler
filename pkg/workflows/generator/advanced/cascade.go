package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func addCascadeFailures(a *Generator, targetsList []targets.Target, r *rand.Rand, params phaseParams) []generator.Stage {
	phaseFailures := make([]Failure, len(a.failures))
	copy(phaseFailures, a.failures)

	stages := make([]generator.Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageTargets := make([]targets.Target, len(targetsList))
		copy(stageTargets, targetsList)

		actions := make([]generator.Action, 0)
		points := a.budget.MaxPoints

		failure := popRandomFailure(phaseFailures, r)
		cost := calculateCost(a.modifiers, failure)

		for i := 0; i < a.retries; i++ {
			if cost*2 <= points {
				break
			}

			failure := popRandomFailure(phaseFailures, r)
			cost = calculateCost(a.modifiers, failure)
		}

		for len(actions) < a.budget.MaxFailures {
			if len(stageTargets) == 0 {
				break
			}

			target := popRandomTarget(stageTargets, r)

			actions = append(actions, generator.Action{
				Info:   failure.Preset.Info(),
				Target: target,
				Engine: failure.Preset.Engine(target, params.StageDuration),
			})

			points -= cost
			if cost > points {
				break
			}
		}

		stages = append(stages, generator.Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
