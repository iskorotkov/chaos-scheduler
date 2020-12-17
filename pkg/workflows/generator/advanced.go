package generator

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Scale int

type Severity int

type Cost float32

const (
	ScaleContainer Scale = iota
	ScalePod
	ScaleDeploymentPart
	ScaleDeployment
	ScaleNode
)

const (
	SeverityNonCritical Severity = iota
	SeverityCritical
	SeverityLethal
)

type Budget struct {
	MaxExperiments int
	MaxPoints      Cost
}

type Modifiers struct {
	ByScale    map[Scale]Cost
	BySeverity map[Severity]Cost
}

type Failure interface {
	Info() experiments.Info
	Scale() Scale
	Severity() Severity
	Engine(t targets.Target, d time.Duration) experiments.Engine
}

type AdvancedGenerator struct {
	stageDuration time.Duration
	retries       int
	failures      []Failure
	seeker        targets.Seeker
	budget        Budget
	modifiers     Modifiers
	logger        *zap.SugaredLogger
}

type PhaseParams struct {
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

	stages := make([]Stage, 0)

	stages = append(stages, a.addIsolatedFailures(t, r, params)...)
	stages = append(stages, a.addCascadeFailures(t, r, params)...)
	stages = append(stages, a.addComplexFailures(t, r, params)...)

	return Scenario{Stages: stages}, nil
}

func (a AdvancedGenerator) addIsolatedFailures(t []targets.Target, r *rand.Rand, params Params) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := a.pickRandomFailure(r)
		target := a.selectTarget(t, r)

		actions := []Action{{
			Info:   failure.Info(),
			Target: target,
			Engine: failure.Engine(target, a.stageDuration),
		}}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: a.stageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) addCascadeFailures(t []targets.Target, r *rand.Rand, params Params) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		actions := make([]Action, 0)

		points := a.budget.MaxPoints
		failure := a.pickRandomFailure(r)

		for i := 0; i < a.retries; i++ {
			cost := a.calculateCost(failure)
			if cost <= points {
				points -= cost
				break
			}

			failure = a.pickRandomFailure(r)
		}

		for len(actions) <= a.budget.MaxExperiments {
			target := a.selectTarget(t, r)

			actions = append(actions, Action{
				Info:   failure.Info(),
				Target: target,
				Engine: failure.Engine(target, a.stageDuration),
			})
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: a.stageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) addComplexFailures(t []targets.Target, r *rand.Rand, params Params) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		actions := make([]Action, 0)

		points := a.budget.MaxPoints
		retries := a.retries
		for len(actions) <= a.budget.MaxExperiments {
			failure := a.pickRandomFailure(r)
			target := a.selectTarget(t, r)

			cost := a.calculateCost(failure)

			if cost <= points {
				points -= cost

				actions = append(actions, Action{
					Info:   failure.Info(),
					Target: target,
					Engine: failure.Engine(target, a.stageDuration),
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
				Info:   failure.Info(),
				Target: target,
				Engine: failure.Engine(target, a.stageDuration),
			})
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: a.stageDuration,
		})
	}

	return stages
}

func (a AdvancedGenerator) pickRandomFailure(r *rand.Rand) Failure {
	i := r.Intn(len(a.failures))
	return a.failures[i]
}

func (a AdvancedGenerator) calculateCost(f Failure) Cost {
	severity, ok := a.modifiers.BySeverity[f.Severity()]
	if !ok {
		severity = 1
	}

	scale, ok := a.modifiers.ByScale[f.Scale()]
	if !ok {
		scale = 1
	}

	return severity * scale
}

func (a AdvancedGenerator) selectTarget(t []targets.Target, r *rand.Rand) targets.Target {
	i := r.Intn(len(t))
	return t[i]
}
