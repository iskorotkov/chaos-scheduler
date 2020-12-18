package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/pod"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"net/http"
)

func createWorkflowFromForm(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, form, error) {
	f, err := parseScenarioParams(r, logger.Named("params"))
	if err != nil {
		return templates.Workflow{}, form{}, err
	}

	wf, err := generateWorkflow(f, cfg, logger.Named("workflow"))
	if err != nil {
		return templates.Workflow{}, form{}, err
	}

	return wf, f, nil
}

func generateWorkflow(params form, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, error) {
	seeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		return templates.Workflow{}, scenarioGeneratorError
	}

	f := failures(cfg)
	g := generator.NewAdvancedGenerator(f, seeker, logger.Named("generator"))
	s, err := g.Generate(generator.Params{Stages: params.Stages, Seed: params.Seed, StageDuration: cfg.StageDuration})
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"config", cfg,
			"failures", f)

		if err == generator.NonPositiveStagesError || err == generator.TooManyStagesError {
			return templates.Workflow{}, scenarioParamsError
		} else {
			return templates.Workflow{}, scenarioGeneratorError
		}
	}

	ext := assemblerExtensions(cfg, logger)
	a := assemblers.NewModularAssembler(ext)
	wf, err := a.Assemble(s)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", ext)
		return templates.Workflow{}, scenarioGeneratorError
	}

	return wf, nil
}

func assemblerExtensions(cfg *config.Config, logger *zap.SugaredLogger) extensions.List {
	return extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval, logger.Named("monitor"))},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}
}

func failures(cfg *config.Config) []generator.Failure {
	return []generator.Failure{
		{
			Preset:   container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 3000},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 10},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 90},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 10},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 90},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 10},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityNonCritical,
		},
		{
			Preset:   container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 90},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 1},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryConsumption: 1000},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   container.DiskFill{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, FillPercentage: 90},
			Scale:    generator.ScaleContainer,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   pod.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    generator.ScalePod,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   pod.Delete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
			Scale:    generator.ScalePod,
			Severity: generator.SeverityLethal,
		},
		{
			Preset:   node.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 2},
			Scale:    generator.ScaleNode,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   node.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryPercentage: 90},
			Scale:    generator.ScaleNode,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   node.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    generator.ScaleNode,
			Severity: generator.SeverityCritical,
		},
		{
			Preset:   node.Restart{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS},
			Scale:    generator.ScaleNode,
			Severity: generator.SeverityLethal,
		},
	}
}
