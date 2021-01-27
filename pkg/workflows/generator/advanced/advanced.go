package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
)

type Generator struct {
	retries   int
	failures  []failures.Failure
	seeker    targets.Seeker
	budget    Budget
	modifiers Modifiers
	logger    *zap.SugaredLogger
}

type Option func(a *Generator)

func NewGenerator(failures []failures.Failure, seeker targets.Seeker, logger *zap.SugaredLogger, opts ...Option) (*Generator, error) {
	a := &Generator{
		failures: failures,
		seeker:   seeker,
		logger:   logger,
	}

	WithRetries(3)(a)

	WithBudget(Budget{
		MaxFailures: 3,
		MaxPoints:   12,
	})(a)

	WithModifiers(Modifiers{
		ByScale: map[metadata.Scale]Cost{
			metadata.ScaleContainer:      1,
			metadata.ScalePod:            1,
			metadata.ScaleDeploymentPart: 1.5,
			metadata.ScaleDeployment:     2,
			metadata.ScaleNode:           4,
		},
		BySeverity: map[metadata.Severity]Cost{
			metadata.SeverityLight:    1,
			metadata.SeveritySevere:   1.5,
			metadata.SeverityCritical: 2,
		},
	})(a)

	for _, opt := range opts {
		opt(a)
	}

	if len(failures) == 0 {
		return nil, generator.ZeroFailures
	}

	if a.budget.MaxFailures < 1 {
		return nil, generator.MaxFailuresError
	}

	if a.budget.MaxPoints < 1 {
		return nil, generator.MaxPointsError
	}

	if a.retries < 0 {
		return nil, generator.RetriesError
	}

	return a, nil
}

func WithRetries(retries int) Option {
	return func(a *Generator) {
		a.retries = retries
	}
}

func WithModifiers(modifiers Modifiers) Option {
	return func(a *Generator) {
		a.modifiers = modifiers
	}
}

func WithBudget(budget Budget) Option {
	return func(a *Generator) {
		a.budget = budget
	}
}

type phaseParams struct {
	Stages        int
	StageDuration time.Duration
}

func (p phaseParams) Generate(rand *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(phaseParams{
		Stages:        1 + rand.Intn(100),
		StageDuration: time.Duration(1+rand.Int63n(120)) * time.Second,
	})
}

func (a *Generator) Generate(params generator.Params) (generator.Scenario, error) {
	seed, stages, stageDuration := params.Seed, params.Stages, params.StageDuration

	r := rand.New(rand.NewSource(seed))

	if stages <= 0 {
		return generator.Scenario{}, generator.NonPositiveStagesError
	}

	if stages > 100 {
		return generator.Scenario{}, generator.TooManyStagesError
	}

	t, err := a.seeker.Targets()
	if err != nil {
		a.logger.Error(err)
		return generator.Scenario{}, generator.TargetsError
	}

	if len(t) == 0 {
		return generator.Scenario{}, generator.ZeroTargetsError
	}

	isolatedFailures := a.addIsolatedFailures(t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	cascadeFailures := a.addCascadeFailures(t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	complexFailures := a.addComplexFailures(t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	s := make([]generator.Stage, 0)
	s = append(s, isolatedFailures...)
	s = append(s, cascadeFailures...)
	s = append(s, complexFailures...)

	return generator.Scenario{Stages: s}, nil
}
