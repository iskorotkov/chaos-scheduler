package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func randomTarget(targets []targets.Target, r *rand.Rand) targets.Target {
	targetIndex := r.Intn(len(targets))
	target := targets[targetIndex]
	return target
}

func randomFailure(failures []experiments.Failure, r *rand.Rand) experiments.Failure {
	failureIndex := r.Intn(len(failures))
	failure := failures[failureIndex]
	return failure
}
