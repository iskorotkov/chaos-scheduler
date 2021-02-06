// Package generate handles generation of chaos scenarios.
package generate

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
)

var (
	ErrNonPositiveStages = errors.New("number of stages must be positive")
	ErrTooManyStages     = errors.New("number of stages can't be that high")
	ErrZeroFailures      = errors.New("can't create scenario out of 0 failures")
	ErrZeroTargets       = errors.New("no targets available")
	ErrMaxFailures       = errors.New("max number of failures must be positive")
	ErrMaxPoints         = errors.New("max points per stage must be positive")
	ErrStageDuration     = errors.New("stage duration must be at least 1s")
)

const (
	retries = 3
)

type Params struct {
	// Seed to use for selecting both targets and failures.
	Seed int64
	// Stages is a total number of stages.
	Stages        int
	StageDuration time.Duration
	Failures      []failures.Failure
	Targets       []targets.Target
	// Budget is a set of restrictions on scenario max damage.
	Budget Budget
	// Modifiers to calculate chaos score of failures.
	Modifiers Modifiers
	Logger    *zap.SugaredLogger
}

func (p Params) Generate(r *rand.Rand, size int) reflect.Value {
	var fs []failures.Failure
	for i := 0; i < r.Intn(10); i++ {
		fs = append(fs, failures.Failure{}.Generate(r, size).Interface().(failures.Failure))
	}

	var ts []targets.Target
	for i := 0; i < r.Intn(10); i++ {
		ts = append(ts, targets.Target{}.Generate(r, size).Interface().(targets.Target))
	}

	return reflect.ValueOf(Params{
		Stages:        -10 + r.Intn(120),
		Seed:          0,
		StageDuration: time.Duration(-10+r.Intn(200)) * time.Second,
		Failures:      fs,
		Targets:       ts,
		Budget:        DefaultBudget(),
		Modifiers:     DefaultModifiers(),
		Logger:        zap.NewNop().Sugar(),
	})
}

func Generate(params Params) (Scenario, error) {
	if len(params.Failures) == 0 {
		return Scenario{}, ErrZeroFailures
	}

	if params.Budget.MaxFailures < 1 {
		return Scenario{}, ErrMaxFailures
	}

	if params.Budget.MaxPoints < 1 {
		return Scenario{}, ErrMaxPoints
	}

	if params.Stages <= 0 {
		return Scenario{}, ErrNonPositiveStages
	}

	if params.Stages > 100 {
		return Scenario{}, ErrTooManyStages
	}

	if len(params.Targets) == 0 {
		return Scenario{}, ErrZeroTargets
	}

	if params.StageDuration < time.Second {
		return Scenario{}, ErrStageDuration
	}

	rng := rand.New(rand.NewSource(params.Seed))

	s := make([]Stage, 0)
	s = append(s, addIsolatedFailures(params, rng)...)
	s = append(s, addCascadeFailures(params, rng)...)
	s = append(s, addComplexFailures(params, rng)...)

	return Scenario{Stages: s}, nil
}

func randomTarget(targets []targets.Target, r *rand.Rand) targets.Target {
	targetIndex := r.Intn(len(targets))
	target := targets[targetIndex]
	return target
}

func randomFailure(failures []failures.Failure, r *rand.Rand) failures.Failure {
	failureIndex := r.Intn(len(failures))
	failure := failures[failureIndex]
	return failure
}
