package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func (a *Generator) addCascadeFailures(targetsList []targets.Target, r *rand.Rand, params phaseParams) []generator.Stage {
	phaseFailures := a.failures

	stages := make([]generator.Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageTargets := targetsList

		actions := make([]generator.Action, 0)
		points := a.budget.MaxPoints

		failure := randomFailure(phaseFailures, r)
		cost := calculateCost(a.modifiers, failure)

		for i := 0; i < a.retries; i++ {
			if cost*2 <= points {
				break
			}

			failure := randomFailure(phaseFailures, r)
			cost = calculateCost(a.modifiers, failure)
		}

		for len(actions) < a.budget.MaxFailures {
			if len(stageTargets) == 0 {
				break
			}

			target := randomTarget(stageTargets, r)

			actions = append(actions, generator.Action{
				Name:     failure.Name(),
				Severity: failure.Severity,
				Scale:    failure.Scale,
				Target:   target,
				Engine:   failure.Template.Instantiate(target, params.StageDuration),
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
