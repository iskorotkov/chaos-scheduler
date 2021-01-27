package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
)

var (
	ScenarioParamsError    = errors.New("couldn't create scenario with given parameters")
	TargetsFetchError      = errors.New("couldn't fetch targets")
	NotEnoughTargetsError  = errors.New("not enough targets present")
	NotEnoughFailuresError = errors.New("not enough failures provided to scenario generator")
	AssembleError          = errors.New("couldn't generate scenario due to unknown reason")
)

//goland:noinspection GoUnusedExportedFunction

//goland:noinspection GoUnusedExportedFunction

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
	ts, err := targets.List(params.AppNS, params.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error())
		if err == targets.ClientsetError {
			return generate.Scenario{}, TargetsFetchError
		} else {
			return generate.Scenario{}, NotEnoughTargetsError
		}
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
			return generate.Scenario{}, NotEnoughTargetsError
		}

		if err == generate.ZeroFailures {
			return generate.Scenario{}, NotEnoughFailuresError
		}

		return generate.Scenario{}, ScenarioParamsError
	}
	return scenario, nil
}

func CreateWorkflow(sp ScenarioParams, wp WorkflowParams, logger *zap.SugaredLogger) (assemble.Workflow, error) {
	scenario, err := CreateScenario(sp, logger)
	if err != nil {
		return assemble.Workflow{}, err
	}

	workflow, err := assemble.Assemble(scenario, wp.Extensions)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", wp.Extensions)
		return assemble.Workflow{}, AssembleError
	}

	return workflow, nil
}
