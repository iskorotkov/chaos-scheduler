package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
)

var (
	ScenarioParamsError    = errors.New("couldn't create scenario with given parameters")
	TargetsSeekerError     = errors.New("couldn't create targets seeker")
	TargetsError           = errors.New("not enough targets present")
	FailuresError          = errors.New("not enough failures provided to scenario generator")
	UnknownGenerationError = errors.New("couldn't generate scenario due to unknown reason")
)

type ScenarioParams struct {
	Seed          int64
	Stages        int
	AppNS         string
	AppLabel      string
	StageDuration time.Duration
	Failures      []failures.Failure
}

func (s ScenarioParams) Generate(rand *rand.Rand, size int) reflect.Value {
	var fs []failures.Failure
	for i := 0; i <= rand.Intn(10); i++ {
		fs = append(fs, failures.Failure{}.Generate(rand, size).Interface().(failures.Failure))
	}

	return reflect.ValueOf(ScenarioParams{
		Seed:          rand.Int63(),
		Stages:        -10 + rand.Intn(120),
		AppNS:         "chaos-app",
		AppLabel:      "app",
		StageDuration: time.Duration(-10+rand.Int63n(200)) * time.Second,
		Failures:      fs,
	})
}

type WorkflowParams struct {
	Extensions assemble.Extensions
}

func (w WorkflowParams) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(WorkflowParams{
		Extensions: assemble.Extensions{}.Generate(rand, size).Interface().(assemble.Extensions),
	})
}

func CreateScenario(params ScenarioParams, logger *zap.SugaredLogger) (generate.Scenario, error) {
	targetSeeker, err := targets.NewSeeker(params.AppNS, params.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error())
		return generate.Scenario{}, TargetsSeekerError
	}

	ts, err := targetSeeker.Targets()
	if err != nil {
		logger.Errorw(err.Error())
		return generate.Scenario{}, TargetsError
	}

	scenario, err := generate.Generate(generate.Params{
		RNG:           rand.New(rand.NewSource(params.Seed)),
		Stages:        params.Stages,
		StageDuration: params.StageDuration,
		Failures:      params.Failures,
		Targets:       ts,
		Retries:       generate.DefaultRetries(),
		Budget:        generate.DefaultBudget(),
		Modifiers:     generate.DefaultModifiers(),
		Logger:        logger.Named("generate"),
	})
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"failures", params.Failures)

		if err == generate.ZeroTargetsError {
			return generate.Scenario{}, TargetsError
		}

		if err == generate.ZeroFailures {
			return generate.Scenario{}, FailuresError
		}

		return generate.Scenario{}, ScenarioParamsError
	}
	return scenario, nil
}

func CreateWorkflow(sp ScenarioParams, wp WorkflowParams, logger *zap.SugaredLogger) (templates.Workflow, error) {
	scenario, err := CreateScenario(sp, logger)
	if err != nil {
		return templates.Workflow{}, err
	}

	workflow, err := assemble.Assemble(scenario, wp.Extensions)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", wp.Extensions)
		return templates.Workflow{}, UnknownGenerationError
	}

	return workflow, nil
}
