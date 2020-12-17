package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
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
	presetList := generator.PresetsList{
		ContainerPresets: []experiments.ContainerPreset{
			container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 100},
			container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 100},
			container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 100},
			container.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 1},
			container.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryConsumption: 1000},
			container.DiskFill{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, FillPercentage: 100},
		},
		PodPresets: []experiments.PodPreset{
			pod.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 100},
			pod.Delete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
		},
		NodePreset: []experiments.NodePreset{
			node.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 2},
			node.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryPercentage: 90},
			node.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 100},
			node.Restart{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS},
		},
	}

	extensionsList := extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval, logger.Named("monitor"))},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}

	seeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, logger.Named("targets"))
	if err != nil {
		logger.Errorw(err.Error(),
			"config", cfg)
		return templates.Workflow{}, scenarioGeneratorError
	}

	g := generator.NewRoundRobin(presetList, seeker, logger.Named("generator"))
	s, err := g.Generate(generator.Params{Stages: params.Stages, Seed: params.Seed, StageDuration: cfg.StageDuration})
	if err != nil {
		logger.Errorw(err.Error(),
			"params", params,
			"config", cfg,
			"presets", presetList)

		if err == generator.NonPositiveStagesError || err == generator.TooManyStagesError {
			return templates.Workflow{}, scenarioParamsError
		} else {
			return templates.Workflow{}, scenarioGeneratorError
		}
	}

	a := assemblers.NewModularAssembler(extensionsList)
	wf, err := a.Assemble(s)
	if err != nil {
		logger.Errorw(err.Error(),
			"extensions", extensionsList)
		return templates.Workflow{}, scenarioGeneratorError
	}

	return wf, nil
}
