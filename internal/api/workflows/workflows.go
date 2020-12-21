package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/pod"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator/advanced"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"net/http"
)

type generatedWorkflow struct {
	Workflow templates.Workflow `json:"workflow"`
	Scenario generator.Scenario `json:"scenario"`
	Params   workflowParams     `json:"params"`
}

func generateWorkflow(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (generatedWorkflow, error) {
	workflowParams, err := parseWorkflowParams(r, logger.Named("params"))
	if err != nil {
		return generatedWorkflow{}, err
	}

	targetSeeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		return generatedWorkflow{}, targetsSeekerError
	}

	failures := enabledFailures(cfg)

	scenarioGenerator, err := advanced.NewGenerator(failures, targetSeeker, logger.Named("generator"))
	if err != nil {
		logger.Errorw(err.Error(),
			"failures", failures)
		return generatedWorkflow{}, scenarioParamsError
	}

	scenario, err := scenarioGenerator.Generate(workflowParams.Stages, workflowParams.Seed, cfg.StageDuration)
	if err != nil {
		logger.Errorw(err.Error(),
			"params", workflowParams,
			"config", cfg,
			"failures", failures)
		return generatedWorkflow{}, scenarioParamsError
	}

	ext := enabledExtensions(cfg, logger)
	a := assemblers.NewModularAssembler(ext)

	workflow, err := a.Assemble(scenario)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", ext)
		return generatedWorkflow{}, workflowGenerationError
	}

	return generatedWorkflow{
		Workflow: workflow,
		Scenario: scenario,
		Params:   workflowParams,
	}, nil
}

func enabledExtensions(cfg *config.Config, logger *zap.SugaredLogger) extensions.List {
	return extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval, logger.Named("monitor"))},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}
}

func enabledFailures(cfg *config.Config) []advanced.Failure {
	return []advanced.Failure{
		{
			Preset:   container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 3000},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 10},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 90},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 10},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 90},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 10},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 90},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 1},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryConsumption: 1000},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   container.DiskFill{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, FillPercentage: 90},
			Scale:    advanced.ScaleContainer,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   pod.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    advanced.ScalePod,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   pod.Delete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
			Scale:    advanced.ScalePod,
			Severity: advanced.SeverityLethal,
		},
		{
			Preset:   node.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 2},
			Scale:    advanced.ScaleNode,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   node.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryPercentage: 90},
			Scale:    advanced.ScaleNode,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   node.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    advanced.ScaleNode,
			Severity: advanced.SeverityCritical,
		},
		{
			Preset:   node.Restart{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS},
			Scale:    advanced.ScaleNode,
			Severity: advanced.SeverityLethal,
		},
	}
}
