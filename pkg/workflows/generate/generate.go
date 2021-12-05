// Package generate handles generation of chaos scenarios.
package generate

import (
	"errors"
	"math/rand"
	"reflect"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
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

type Seeds struct {
	Targets  int64 `json:"targets"`
	Failures int64 `json:"failures"`
}

type Stages struct {
	Single  int `json:"single"`
	Similar int `json:"similar"`
	Mixed   int `json:"mixed"`
}

func (s Stages) Sum() int {
	return s.Single + s.Similar + s.Mixed
}

type Params struct {
	// Seed to use for selecting both targets and failures.
	Seed Seeds
	// Stages is a total number of stages.
	Stages        Stages
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
		Stages: Stages{
			Single:  -10 + r.Intn(120),
			Similar: -10 + r.Intn(120),
			Mixed:   -10 + r.Intn(120),
		},
		Seed: Seeds{
			Targets:  0,
			Failures: 0,
		},
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

	if params.Stages.Single < 0 || params.Stages.Similar < 0 || params.Stages.Mixed < 0 || params.Stages.Sum() <= 0 {
		return Scenario{}, ErrNonPositiveStages
	}

	if params.Stages.Single > 100 || params.Stages.Similar > 100 || params.Stages.Mixed > 100 {
		return Scenario{}, ErrTooManyStages
	}

	if len(params.Targets) == 0 {
		return Scenario{}, ErrZeroTargets
	}

	if params.StageDuration < time.Second {
		return Scenario{}, ErrStageDuration
	}

	failuresRng := rand.New(rand.NewSource(params.Seed.Failures))
	targetsRng := rand.New(rand.NewSource(params.Seed.Targets))

	s := make([]Stage, 0)
	s = append(s, addIsolatedFailures(params, failuresRng, targetsRng)...)
	s = append(s, addCascadeFailures(params, failuresRng, targetsRng)...)
	s = append(s, addComplexFailures(params, failuresRng, targetsRng)...)

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
