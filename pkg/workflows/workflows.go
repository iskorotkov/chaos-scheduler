package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator/advanced"
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
	Extensions extensions.Extensions
}

func (w WorkflowParams) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(WorkflowParams{
		Extensions: extensions.Extensions{}.Generate(rand, size).Interface().(extensions.Extensions),
	})
}

func CreateScenario(params ScenarioParams, logger *zap.SugaredLogger) (generator.Scenario, error) {
	targetSeeker, err := targets.NewSeeker(params.AppNS, params.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error())
		return generator.Scenario{}, TargetsSeekerError
	}

	scenarioGenerator, err := advanced.NewGenerator(params.Failures, targetSeeker, logger.Named("generator"))
	if err != nil {
		logger.Errorw(err.Error(),
			"failures", params.Failures)
		return generator.Scenario{}, ScenarioParamsError
	}

	scenario, err := scenarioGenerator.Generate(params.Stages, params.Seed, params.StageDuration)
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"failures", params.Failures)

		if err == generator.ZeroTargetsError {
			return generator.Scenario{}, TargetsError
		}

		if err == generator.ZeroFailures {
			return generator.Scenario{}, FailuresError
		}

		return generator.Scenario{}, ScenarioParamsError
	}
	return scenario, nil
}

func CreateWorkflow(sp ScenarioParams, wp WorkflowParams, logger *zap.SugaredLogger) (templates.Workflow, error) {
	scenario, err := CreateScenario(sp, logger)
	if err != nil {
		return templates.Workflow{}, err
	}

	a := assemblers.NewModularAssembler(wp.Extensions, logger.Named("assembler"))

	workflow, err := a.Assemble(scenario)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", wp.Extensions)
		return templates.Workflow{}, UnknownGenerationError
	}

	return workflow, nil
}
