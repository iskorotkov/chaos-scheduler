package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/concrete"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

func generateWorkflow(form Form, cfg *config.Config) (templates.Workflow, error) {
	presetList := experiments.List{
		PodPresets: []experiments.PodEnginePreset{
			concrete.PodDelete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
		},
		ContainerPresets: []experiments.ContainerEnginePreset{
			concrete.PodNetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			concrete.PodNetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 100},
		},
	}

	extensionsList := extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval)},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}

	seeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, cfg.IsInKubernetes)
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, ScenarioGeneratorError
	}

	g := generators.NewRoundRobinGenerator(presetList, seeker)
	s, err := g.Generate(generators.Params{Stages: form.Stages, Seed: form.Seed, StageDuration: cfg.StageDuration})
	if err != nil {
		logger.Error(err)

		if err == generators.NonPositiveStagesError || err == generators.TooManyStagesError {
			return templates.Workflow{}, ScenarioParamsError
		}

		return templates.Workflow{}, ScenarioGeneratorError
	}

	a := assemblers.NewModularAssembler(extensionsList)
	wf, err := a.Assemble(s)
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, ScenarioGeneratorError
	}

	return wf, nil
}
