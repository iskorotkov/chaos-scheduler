package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Generator struct {
	retries   int
	failures  []experiments.Failure
	seeker    targets.Seeker
	budget    Budget
	modifiers Modifiers
	logger    *zap.SugaredLogger
}

type Option func(a *Generator)

func NewGenerator(failures []experiments.Failure, seeker targets.Seeker, logger *zap.SugaredLogger, opts ...Option) (*Generator, error) {
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
			metadata.SeverityNonCritical: 1,
			metadata.SeverityCritical:    1.5,
			metadata.SeverityLethal:      2,
		},
	})(a)

	for _, opt := range opts {
		opt(a)
	}

	if len(failures) == 0 {
		return nil, ZeroFailures
	}

	if a.budget.MaxFailures < 1 {
		return nil, MaxFailuresError
	}

	if a.budget.MaxPoints < 1 {
		return nil, MaxPointsError
	}

	if a.retries < 0 {
		return nil, RetriesError
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

func (a *Generator) Generate(stages int, seed int64, stageDuration time.Duration) (generator.Scenario, error) {
	r := rand.New(rand.NewSource(seed))
	t, err := a.seeker.Targets()

	if stages <= 0 {
		return generator.Scenario{}, NonPositiveStagesError
	}

	if stages > 100 {
		return generator.Scenario{}, TooManyStagesError
	}

	if err != nil {
		a.logger.Error(err)
		return generator.Scenario{}, TargetsError
	}

	if len(t) == 0 {
		return generator.Scenario{}, LowTargetsError
	}

	isolatedFailures := addIsolatedFailures(a, t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	cascadeFailures := addCascadeFailures(a, t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	complexFailures := addComplexFailures(a, t, r, phaseParams{
		Stages:        stages,
		StageDuration: stageDuration,
	})

	s := make([]generator.Stage, 0)
	s = append(s, isolatedFailures...)
	s = append(s, cascadeFailures...)
	s = append(s, complexFailures...)

	return generator.Scenario{Stages: s}, nil
}
