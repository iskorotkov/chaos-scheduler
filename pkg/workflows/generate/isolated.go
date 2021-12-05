package generate

import (
	"math/rand"
)

// addIsolatedFailures adds one failure in each stage.
func addIsolatedFailures(params Params, rng *rand.Rand) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := randomFailure(params.Failures, rng)
		target := randomTarget(params.Targets, rng)

		actions := []Action{{
			Name:     failure.Blueprint.Name(),
			Severity: failure.Severity,
			Scale:    failure.Scale,
			Target:   target,
			Engine:   failure.Blueprint.Instantiate(target, params.StageDuration),
		}}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
