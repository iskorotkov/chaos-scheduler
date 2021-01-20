package workflows

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/pod"
	"go.uber.org/zap"
)

func enabledExtensions(cfg *config.Config, logger *zap.SugaredLogger) extensions.List {
	return extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval, logger.Named("monitor"))},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}
}

func enabledFailures(cfg *config.Config) []failures.Failure {
	return []failures.Failure{
		{
			Template: container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityNonCritical,
		},
		{
			Template: container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 3000},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 10},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityNonCritical,
		},
		{
			Template: container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 90},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 10},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityNonCritical,
		},
		{
			Template: container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 90},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 10},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityNonCritical,
		},
		{
			Template: container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 90},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 1},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryConsumption: 1000},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: container.DiskFill{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, FillPercentage: 90},
			Scale:    metadata.ScaleContainer,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: pod.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    metadata.ScalePod,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: pod.Delete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
			Scale:    metadata.ScalePod,
			Severity: metadata.SeverityLethal,
		},
		{
			Template: node.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 2},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: node.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryPercentage: 90},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: node.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeverityCritical,
		},
		{
			Template: node.Restart{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeverityLethal,
		},
	}
}
