package handlers

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/container"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/node"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/pod"
	"go.uber.org/zap"
)

// enabledExtensions are assembling extension used by default.
func enabledExtensions(cfg *config.Config, logger *zap.SugaredLogger) assemble.ExtCollection {
	return assemble.ExtCollection{
		Action:   nil,
		Stage:    []assemble.StageExt{assemble.UseMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.AppLabel, cfg.StageInterval, logger.Named("monitor"))},
		Workflow: []assemble.WorkflowExt{assemble.UseSteps()},
	}
}

// enabledFailures are chaos failures used by default.
func enabledFailures(cfg *config.Config) []failures.Failure {
	return []failures.Failure{
		{
			Blueprint: container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeverityLight,
		},
		{
			Blueprint: container.NetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 3000},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 10},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeverityLight,
		},
		{
			Blueprint: container.NetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 90},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 10},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeverityLight,
		},
		{
			Blueprint: container.NetworkCorruption{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, CorruptionPercentage: 90},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 10},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeverityLight,
		},
		{
			Blueprint: container.NetworkDuplication{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, DuplicationPercentage: 90},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: container.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 1},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: container.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryConsumption: 1000},
			Scale:     metadata.ScaleContainer,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: pod.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:     metadata.ScalePod,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: pod.Delete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
			Scale:     metadata.ScalePod,
			Severity:  metadata.SeverityCritical,
		},
		{
			Blueprint: node.CPUHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Cores: 2},
			Scale:     metadata.ScaleNode,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: node.MemoryHog{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, MemoryPercentage: 90},
			Scale:     metadata.ScaleNode,
			Severity:  metadata.SeveritySevere,
		},
		{
			Blueprint: node.IOStress{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, UtilizationPercentage: 90},
			Scale:     metadata.ScaleNode,
			Severity:  metadata.SeveritySevere,
		},
	}
}
