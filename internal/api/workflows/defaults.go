package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/pod"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator/advanced"
	"go.uber.org/zap"
)

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
