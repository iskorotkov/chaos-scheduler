package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines/factories"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/exporters"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/scenarios"
)

var (
	TemplatesImportError    = errors.New("couldn't import templates")
	ScenarioGenerationError = errors.New("couldn't generate scenario")
	WorkflowAssemblingError = errors.New("couldn't assemble workflow")
	WorkflowExportError     = errors.New("couldn't export workflow")
)

type WorkflowParams struct {
	Generator scenarios.Generator
	Config    scenarios.ScenarioParams
	Assembler assemblers.Assembler
	Exporter  exporters.Exporter
}

func NewWorkflow(params WorkflowParams) (string, error) {
	fs := []engines.Factory{
		factories.PodDeleteFactory{Namespace: "litmus", TargetNamespace: "chaos-app", Duration: 60, Interval: 5, Force: false},
	}

	scenarioConfig := scenarios.ScenarioParams{Stages: params.Config.Stages, Seed: params.Config.Seed}
	scenario, err := params.Generator.Generate(fs, scenarioConfig)
	if err != nil {
		logger.Error(err)
		return "", ScenarioGenerationError
	}

	workflow, err := params.Assembler.Assemble(scenario)
	if err != nil {
		logger.Error(err)
		return "", WorkflowAssemblingError
	}

	str, err := params.Exporter.Export(workflow)
	if err != nil {
		logger.Error(err)
		return "", WorkflowExportError
	}

	return str, nil
}
