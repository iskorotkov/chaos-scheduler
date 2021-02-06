package generate

import (
	"math/rand"
)

// addComplexFailures add several failures of different types in each stage.
func addComplexFailures(params Params, rng *rand.Rand) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageFailures := params.Failures
		stageTargets := params.Targets

		actions := make([]Action, 0)

		points := params.Budget.MaxPoints
		stageRetries := retries
		for len(actions) < params.Budget.MaxFailures {
			failure := randomFailure(stageFailures, rng)
			target := randomTarget(stageTargets, rng)
			cost := calculateCost(params.Modifiers, failure)

			if cost <= points {
				points -= cost

				actions = append(actions, Action{
					Name:     failure.Name(),
					Severity: failure.Severity,
					Scale:    failure.Scale,
					Target:   target,
					Engine:   failure.Blueprint.Instantiate(target, params.StageDuration),
				})
			} else {
				if stageRetries <= 0 {
					break
				}

				stageRetries--
			}
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
