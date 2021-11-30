package generate

import (
	"math/rand"
)

// addIsolatedFailures adds one failure in each stage.
func addIsolatedFailures(params Params, failuresRng *rand.Rand, targetsRng *rand.Rand) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages.Single; i++ {
		failure := randomFailure(params.Failures, failuresRng)
		target := randomTarget(params.Targets, targetsRng)

		actions := []Action{{
			Name:     failure.Blueprint.Name(),
			Type:     failure.Blueprint.Type(),
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
