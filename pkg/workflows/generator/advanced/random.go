package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"math/rand"
)

func popRandomTarget(targets []targets.Target, r *rand.Rand) targets.Target {
	targetIndex := r.Intn(len(targets))
	target := targets[targetIndex]
	targets = append(targets[:targetIndex], targets[targetIndex+1:]...)
	return target
}

func popRandomFailure(failures []Failure, r *rand.Rand) Failure {
	failureIndex := r.Intn(len(failures))
	failure := failures[failureIndex]
	failures = append(failures[:failureIndex], failures[failureIndex+1:]...)
	return failure
}
