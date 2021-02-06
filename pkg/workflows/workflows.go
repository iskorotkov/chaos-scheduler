// Package workflows handles scenario generation and workflow execution.
package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"time"
)

var (
	// ErrScenarioParams is returned when scenario params are invalid.
	ErrScenarioParams = errors.New("couldn't create scenario with given parameters")
	// ErrTargetsFetch is returned when targets couldn't be fetched.
	ErrTargetsFetch = errors.New("couldn't fetch targets")
	// ErrNotEnoughTargets is returned when the number of targets is too small.
	ErrNotEnoughTargets = errors.New("not enough targets present")
	// ErrNotEnoughFailures is returned when the number of failures is too small.
	ErrNotEnoughFailures = errors.New("not enough failures provided to scenario generator")
	// ErrAssemble is returned when Argo workflow preparation from generated scenario has failed.
	ErrAssemble = errors.New("couldn't generate scenario due to unknown reason")
	// ErrExecution is returned when execution of the generated workflow failed.
	ErrExecution = errors.New("couldn't execute generated workflow")
)

// ScenarioParams describes all values required to generate scenario.
type ScenarioParams struct {
	// Seed is used for selecting both targets and failures.
	Seed int64
	// Stages is a number of stages in a generated scenario.
	Stages int
	// AppNS is a namespace with targets.
	AppNS string
	// AppLabel is a label used for target selection.
	AppLabel string
	// StageDuration is a duration of each stage.
	StageDuration time.Duration
	// Failures is a list of enabled failures.
	Failures []failures.Failure
	// TargetFinder is used to fetch targets.
	TargetFinder targets.TargetFinder
}

// Generate returns random ScenarioParams.
func (s ScenarioParams) Generate(rand *rand.Rand, size int) reflect.Value {
	var fs []failures.Failure
	for i := 0; i <= rand.Intn(10); i++ {
		fs = append(fs, failures.Failure{}.Generate(rand, size).Interface().(failures.Failure))
	}

	finder := targets.TestTargetFinder{}.Generate(rand, size).Interface().(targets.TestTargetFinder)
	return reflect.ValueOf(ScenarioParams{
		Seed:          rand.Int63(),
		Stages:        -10 + rand.Intn(120),
		AppNS:         "chaos-app",
		AppLabel:      "app",
		StageDuration: time.Duration(-10+rand.Int63n(200)) * time.Second,
		Failures:      fs,
		TargetFinder:  &finder,
	})
}

// CreateScenario returns new generate.Scenario to preview.
func CreateScenario(params ScenarioParams, logger *zap.SugaredLogger) (generate.Scenario, error) {
	ts, err := params.TargetFinder.List(params.AppNS, params.AppLabel)
	if err != nil {
		logger.Errorw(err.Error())
		if err == targets.ErrClient {
			return generate.Scenario{}, ErrTargetsFetch
		} else {
			return generate.Scenario{}, ErrNotEnoughTargets
		}
	}

	scenario, err := generate.Generate(generate.Params{
		Seed:          params.Seed,
		Stages:        params.Stages,
		StageDuration: params.StageDuration,
		Failures:      params.Failures,
		Targets:       ts,
		Budget:        generate.DefaultBudget(),
		Modifiers:     generate.DefaultModifiers(),
		Logger:        logger.Named("generate"),
	})
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"failures", params.Failures)

		if err == generate.ErrZeroTargets {
			return generate.Scenario{}, ErrNotEnoughTargets
		}

		if err == generate.ErrZeroFailures {
			return generate.Scenario{}, ErrNotEnoughFailures
		}

		return generate.Scenario{}, ErrScenarioParams
	}
	return scenario, nil
}

// WorkflowParams describes all values required to generate workflow from a scenario.
type WorkflowParams struct {
	// Extensions describes enabled assembler extensions.
	Extensions assemble.ExtCollection
}

// Generate returns random WorkflowParams.
func (w WorkflowParams) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(WorkflowParams{
		Extensions: assemble.ExtCollection{}.Generate(rand, size).Interface().(assemble.ExtCollection),
	})
}

// CreateWorkflow returns new assemble.Workflow for previewing or execution.
func CreateWorkflow(sp ScenarioParams, wp WorkflowParams, logger *zap.SugaredLogger) (assemble.Workflow, error) {
	scenario, err := CreateScenario(sp, logger)
	if err != nil {
		return assemble.Workflow{}, err
	}

	workflow, err := assemble.Assemble(scenario, wp.Extensions)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", wp.Extensions)
		return assemble.Workflow{}, ErrAssemble
	}

	return workflow, nil
}

// ExecutionParams describes all values required to execute generated workflow.
type ExecutionParams struct {
	// Executor handles workflow execution.
	Executor execute.Executor
}

// Generate returns random ExecutionParams.
func (e ExecutionParams) Generate(rand *rand.Rand, size int) reflect.Value {
	executor := execute.TestExecutor{}.Generate(rand, size).Interface().(execute.TestExecutor)
	return reflect.ValueOf(ExecutionParams{
		Executor: &executor,
	})
}

// ExecuteWorkflow generates q new workflow and executes it.
func ExecuteWorkflow(sp ScenarioParams, wp WorkflowParams, ep ExecutionParams, logger *zap.SugaredLogger) (assemble.Workflow, error) {
	workflow, err := CreateWorkflow(sp, wp, logger.Named("create-workflow"))
	if err != nil {
		return assemble.Workflow{}, err
	}

	workflow, err = ep.Executor.Execute(workflow)
	if err != nil {
		logger.Errorw(err.Error())
		return assemble.Workflow{}, ErrExecution
	}

	return workflow, err
}
