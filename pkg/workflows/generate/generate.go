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
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroFailures           = errors.New("can't create scenario out of 0 failures")
	ZeroTargetsError       = errors.New("no targets available")
	MaxFailuresError       = errors.New("max number of failures must be positive")
	MaxPointsError         = errors.New("max points per stage must be positive")
	RetriesError           = errors.New("retries must be non negative")
	StageDurationError     = errors.New("stage duration must be at least 1s")
)

func DefaultRetries() int {
	return 3
}

type Params struct {
	RNG           *rand.Rand
	Stages        int
	StageDuration time.Duration
	Failures      []failures.Failure
	Targets       []targets.Target
	Retries       int
	Budget        Budget
	Modifiers     Modifiers
	Logger        *zap.SugaredLogger
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
		RNG:           rand.New(rand.NewSource(r.Int63())),
		StageDuration: time.Duration(-10+r.Intn(200)) * time.Second,
		Failures:      fs,
		Targets:       ts,
		Retries:       -5 + r.Intn(20),
		Budget:        DefaultBudget(),
		Modifiers:     DefaultModifiers(),
		Logger:        zap.NewNop().Sugar(),
	})
}

func Generate(params Params) (Scenario, error) {
	if len(params.Failures) == 0 {
		return Scenario{}, ZeroFailures
	}

	if params.Budget.MaxFailures < 1 {
		return Scenario{}, MaxFailuresError
	}

	if params.Budget.MaxPoints < 1 {
		return Scenario{}, MaxPointsError
	}

	if params.Retries < 0 {
		return Scenario{}, RetriesError
	}

	if params.Stages <= 0 {
		return Scenario{}, NonPositiveStagesError
	}

	if params.Stages > 100 {
		return Scenario{}, TooManyStagesError
	}

	if len(params.Targets) == 0 {
		return Scenario{}, ZeroTargetsError
	}

	if params.StageDuration < time.Second {
		return Scenario{}, StageDurationError
	}

	s := make([]Stage, 0)
	s = append(s, addIsolatedFailures(params)...)
	s = append(s, addCascadeFailures(params)...)
	s = append(s, addComplexFailures(params)...)

	return Scenario{Stages: s}, nil
}
