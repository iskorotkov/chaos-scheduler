package argov2

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/exporters"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/importers"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

var (
	TemplatesImportError    = errors.New("couldn't import templates")
	ScenarioGenerationError = errors.New("couldn't generate scenario")
	WorkflowAssemblingError = errors.New("couldn't assemble workflow")
	WorkflowExportError     = errors.New("couldn't export workflow")
)

type Config struct {
	Importer  importers.Importer
	Generator scenarios.Generator
	Config    scenarios.Config
	Assembler assemblers.Assembler
	Exporter  exporters.Exporter
}

func NewWorkflow(config Config) (string, error) {
	templates, err := config.Importer.Import()
	if err != nil {
		logger.Error(err)
		return "", TemplatesImportError
	}

	scenarioConfig := scenarios.Config{Stages: config.Config.Stages, Seed: config.Config.Seed}
	scenario, err := config.Generator.Generate(templates, scenarioConfig)
	if err != nil {
		logger.Error(err)
		return "", ScenarioGenerationError
	}

	workflow, err := config.Assembler.Assemble(scenario)
	if err != nil {
		logger.Error(err)
		return "", WorkflowAssemblingError
	}

	str, err := config.Exporter.Export(workflow)
	if err != nil {
		logger.Error(err)
		return "", WorkflowExportError
	}

	return str, nil
}
