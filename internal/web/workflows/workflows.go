package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/pod"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator/advanced"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"net/http"
)

type workflow struct {
	Workflow templates.Workflow `json:"workflow"`
	Params   workflowParams     `json:"params"`
}

func createWorkflowFromForm(r *http.Request, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, workflowParams, error) {
	f, err := parseWorkflowParams(r, logger.Named("params"))
	if err != nil {
		return templates.Workflow{}, workflowParams{}, err
	}

	wf, err := generateWorkflow(f, cfg, logger.Named("workflow"))
	if err != nil {
		return templates.Workflow{}, workflowParams{}, err
	}

	return wf, f, nil
}

func generateWorkflow(params workflowParams, cfg *config.Config, logger *zap.SugaredLogger) (templates.Workflow, error) {
	seeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		return templates.Workflow{}, targetsSeekerError
	}

	f := failures(cfg)
	g, err := advanced.NewGenerator(f, seeker, logger.Named("generator"))
	if err != nil {
		logger.Errorw(err.Error(),
			"failures", f)

		return templates.Workflow{}, scenarioParamsError
	}

	s, err := g.Generate(params.Stages, params.Seed, cfg.StageDuration)
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"config", cfg,
			"failures", f)

		return templates.Workflow{}, scenarioParamsError
	}

	ext := assemblerExtensions(cfg, logger)
	a := assemblers.NewModularAssembler(ext)
	wf, err := a.Assemble(s)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", ext)
		return templates.Workflow{}, workflowGenerationError
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

func failures(cfg *config.Config) []advanced.Failure {
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
