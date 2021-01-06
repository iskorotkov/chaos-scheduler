package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator/advanced"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"net/http"
)

func generateWorkflow(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, error) {
	scenario, err := generateScenario(r, cfg, logger)
	if err != nil {
		return templates.Workflow{}, err
	}

	ext := enabledExtensions(cfg, logger)
	a := assemblers.NewModularAssembler(ext)

	workflow, err := a.Assemble(scenario)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", ext)
		return templates.Workflow{}, workflowGenerationError
	}

	return workflow, nil
}

func generateScenario(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (generator.Scenario, error) {
	workflowParams, err := parseWorkflowParams(r, logger.Named("params"))
	if err != nil {
		return generator.Scenario{}, err
	}

	targetSeeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		return generator.Scenario{}, targetsSeekerError
	}

	failures := enabledFailures(cfg)

	scenarioGenerator, err := advanced.NewGenerator(failures, targetSeeker, logger.Named("generator"))
	if err != nil {
		logger.Errorw(err.Error(),
			"failures", failures)
		return generator.Scenario{}, scenarioParamsError
	}

	scenario, err := scenarioGenerator.Generate(workflowParams.Stages, workflowParams.Seed, cfg.StageDuration)
	if err != nil {
		logger.Errorw(err.Error(),
			"params", workflowParams,
			"config", cfg,
			"failures", failures)

		if err == advanced.LowTargetsError {
			return generator.Scenario{}, targetsError
		}

		if err == advanced.ZeroFailures {
			return generator.Scenario{}, failuresError
		}

		return generator.Scenario{}, scenarioParamsError
	}

	return scenario, nil
}
