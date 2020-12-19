package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func addIsolatedFailures(a *Generator, targetsList []targets.Target, r *rand.Rand, params phaseParams) []generator.Stage {
	stages := make([]generator.Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := randomFailure(a.failures, r)
		target := randomTarget(targetsList, r)

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
