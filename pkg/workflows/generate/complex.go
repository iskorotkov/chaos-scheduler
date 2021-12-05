package generate

import (
	"math/rand"
)

// addComplexFailures add several failures of different types in each stage.
func addComplexFailures(params Params, failuresRng *rand.Rand, targetsRng *rand.Rand) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages.Mixed; i++ {
		stageFailures := params.Failures
		stageTargets := params.Targets

		steps := make([]Step, 0)

		points := params.Budget.MaxPoints
		stageRetries := retries
		for len(steps) < params.Budget.MaxFailures {
			failure := randomFailure(stageFailures, failuresRng)
			target := randomTarget(stageTargets, targetsRng)
			cost := calculateCost(params.Modifiers, failure)

			if cost <= points {
				points -= cost

				steps = append(steps, Step{
					Name:     failure.Blueprint.Name(),
					Type:     failure.Blueprint.Type(),
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
			Steps:    steps,
			Duration: params.StageDuration,
		})
	}

	return stages
}
