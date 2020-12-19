package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func addIsolatedFailures(a *Generator, targetsList []targets.Target, r *rand.Rand, params phaseParams) []generator.Stage {
	phaseFailures := make([]Failure, len(a.failures))
	copy(phaseFailures, a.failures)

	stages := make([]generator.Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := popRandomFailure(phaseFailures, r)
		target := popRandomTarget(targetsList, r)

		actions := []generator.Action{{
			Info:   failure.Preset.Info(),
			Target: target,
			Engine: failure.Preset.Engine(target, params.StageDuration),
		}}

		stages = append(stages, generator.Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
