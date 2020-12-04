package workflows

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/exporters"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
)

var (
	TemplatesImportError    = errors.New("couldn't import templates")
	ScenarioGenerationError = errors.New("couldn't generate scenario")
	WorkflowAssemblingError = errors.New("couldn't assemble workflow")
	WorkflowExportError     = errors.New("couldn't export workflow")
)

type WorkflowParams struct {
	Generator generators.Generator
	Assembler assemblers.Assembler
	Exporter  exporters.Exporter
	Params    generators.Params
}

func NewWorkflow(g generators.Generator, a assemblers.Assembler, e exporters.Exporter, params generators.Params) (string, error) {
	scenario, err := g.Generate(params)
	if err != nil {
		logger.Error(err)
		return "", ScenarioGenerationError
	}

	workflow, err := a.Assemble(scenario)
	if err != nil {
		logger.Error(err)
		return "", WorkflowAssemblingError
	}

	str, err := e.Export(workflow)
	if err != nil {
		logger.Error(err)
		return "", WorkflowExportError
	}

	return str, nil
}
