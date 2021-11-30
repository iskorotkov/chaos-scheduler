package generate

import (
	"math/rand"
)

// addCascadeFailures add several failures of the same type in each stage.
func addCascadeFailures(params Params, failuresRng *rand.Rand, targetsRng *rand.Rand) []Stage {
	phaseFailures := params.Failures

	stages := make([]Stage, 0)

	for i := 0; i < params.Stages.Similar; i++ {
		stageTargets := params.Targets

		actions := make([]Action, 0)
		points := params.Budget.MaxPoints

		failure := randomFailure(phaseFailures, failuresRng)
		cost := calculateCost(params.Modifiers, failure)

		for i := 0; i < retries; i++ {
			if cost*2 <= points {
				break
			}

			failure := randomFailure(phaseFailures, failuresRng)
			cost = calculateCost(params.Modifiers, failure)
		}

		for len(actions) < params.Budget.MaxFailures {
			if len(stageTargets) == 0 {
				break
			}

			target := randomTarget(stageTargets, targetsRng)

			actions = append(actions, Action{
				Name:     failure.Blueprint.Name(),
				Type:     failure.Blueprint.Type(),
				Severity: failure.Severity,
				Scale:    failure.Scale,
				Target:   target,
				Engine:   failure.Blueprint.Instantiate(target, params.StageDuration),
			})

			points -= cost
			if cost > points {
				break
			}
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
