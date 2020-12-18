package generator

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type AdvancedGenerator struct {
	retries   int
	failures  []Failure
	seeker    targets.Seeker
	budget    Budget
	modifiers Modifiers
	logger    *zap.SugaredLogger
}

type Option func(a *AdvancedGenerator)

func NewAdvancedGenerator(failures []Failure, seeker targets.Seeker, logger *zap.SugaredLogger, opts ...Option) *AdvancedGenerator {
	a := &AdvancedGenerator{
		failures: failures,
		seeker:   seeker,
		logger:   logger,
	}

	WithRetries(3)(a)

	WithBudget(Budget{
		MaxExperiments: 3,
		MaxPoints:      12,
	})(a)

	WithModifiers(Modifiers{
		ByScale: map[Scale]Cost{
			ScaleContainer:      1,
			ScalePod:            1,
			ScaleDeploymentPart: 1.5,
			ScaleDeployment:     2,
			ScaleNode:           4,
		},
		BySeverity: map[Severity]Cost{
			SeverityNonCritical: 1,
			SeverityCritical:    1.5,
			SeverityLethal:      2,
		},
	})(a)

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func WithRetries(retries int) Option {
	return func(a *AdvancedGenerator) {
		a.retries = retries
	}
}

func WithModifiers(modifiers Modifiers) Option {
	return func(a *AdvancedGenerator) {
		a.modifiers = modifiers
	}
}

func WithBudget(budget Budget) Option {
	return func(a *AdvancedGenerator) {
		a.budget = budget
	}
}

type phaseParams struct {
	Stages        int
	StageDuration time.Duration
}

func (a AdvancedGenerator) Generate(params Params) (Scenario, error) {
	r := rand.New(rand.NewSource(params.Seed))
	t, err := a.seeker.Targets()
	if err != nil {
		a.logger.Error(err)
		return Scenario{}, TargetsError
	}

	isolatedFailures := a.addIsolatedFailures(t, r, phaseParams{
		Stages:        params.Stages,
		StageDuration: params.StageDuration,
	})

	cascadeFailures := a.addCascadeFailures(t, r, phaseParams{
		Stages:        params.Stages,
		StageDuration: params.StageDuration,
	})

	complexFailures := a.addComplexFailures(t, r, phaseParams{
		Stages:        params.Stages,
		StageDuration: params.StageDuration,
	})

	stages := make([]Stage, 0)
	stages = append(stages, isolatedFailures...)
	stages = append(stages, cascadeFailures...)
	stages = append(stages, complexFailures...)

	return Scenario{Stages: stages}, nil
}

func (a AdvancedGenerator) addIsolatedFailures(t []targets.Target, r *rand.Rand, params phaseParams) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := a.pickRandomFailure(r)
		target := a.selectTarget(t, r)

		actions := []Action{{
			Info:   failure.Preset.Info(),
			Target: target,
			Engine: failure.Preset.Engine(target, params.StageDuration),
		}}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) addCascadeFailures(t []targets.Target, r *rand.Rand, params phaseParams) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		actions := make([]Action, 0)

		points := a.budget.MaxPoints
		failure := a.pickRandomFailure(r)
		cost := a.calculateCost(failure)

		for i := 0; i < a.retries; i++ {
			if cost <= points {
				break
			}

			failure = a.pickRandomFailure(r)
			cost = a.calculateCost(failure)
		}

		for len(actions) < a.budget.MaxExperiments {
			target := a.selectTarget(t, r)

			actions = append(actions, Action{
				Info:   failure.Preset.Info(),
				Target: target,
				Engine: failure.Preset.Engine(target, params.StageDuration),
			})

			points -= cost
			if cost > points {
				break
			}
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) addComplexFailures(t []targets.Target, r *rand.Rand, params phaseParams) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		actions := make([]Action, 0)

		points := a.budget.MaxPoints
		retries := a.retries
		for len(actions) < a.budget.MaxExperiments {
			failure := a.pickRandomFailure(r)
			target := a.selectTarget(t, r)

			cost := a.calculateCost(failure)

			if cost <= points {
				points -= cost

				actions = append(actions, Action{
					Info:   failure.Preset.Info(),
					Target: target,
					Engine: failure.Preset.Engine(target, params.StageDuration),
				})
			} else {
				if retries <= 0 {
					break
				}

				retries--
			}
		}

		if len(actions) == 0 {
			failure := a.pickRandomFailure(r)
			target := a.selectTarget(t, r)

			actions = append(actions, Action{
				Info:   failure.Preset.Info(),
				Target: target,
				Engine: failure.Preset.Engine(target, params.StageDuration),
			})
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) pickRandomFailure(r *rand.Rand) Failure {
	i := r.Intn(len(a.failures))
	return a.failures[i]
}

func (a AdvancedGenerator) calculateCost(f Failure) Cost {
	severity, ok := a.modifiers.BySeverity[f.Severity]
	if !ok {
		severity = 1
	}

	scale, ok := a.modifiers.ByScale[f.Scale]
	if !ok {
		scale = 1
	}

	return severity * scale
}

func (a AdvancedGenerator) selectTarget(t []targets.Target, r *rand.Rand) targets.Target {
	i := r.Intn(len(t))
	return t[i]
}
