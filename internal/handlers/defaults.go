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
		Action: nil,
		Stage: []assemble.StageExt{
			assemble.UseMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.AppLabel, cfg.StageInterval, logger.Named("monitor")),
		},
		Workflow: []assemble.WorkflowExt{
			assemble.UseSteps(),
		},
	}
}

// enabledFailures are chaos failures used by default.
func enabledFailures(cfg *config.Config) []failures.Failure {
	podsAffected := []struct {
		Percentage int
		Scale      metadata.Scale
	}{
		{Scale: metadata.ScalePod},
		{Percentage: 50, Scale: metadata.ScaleDeploymentPart},
		{Percentage: 100, Scale: metadata.ScaleDeployment},
	}

	activitiesAffected := []struct {
		Percentage int
		Severity   metadata.Severity
	}{
		{Percentage: 10, Severity: metadata.SeverityLight},
		{Percentage: 90, Severity: metadata.SeveritySevere},
	}

	// Node failures
	fs := []failures.Failure{
		{
			Blueprint: node.CPUHog{
				Namespace:    cfg.ChaosNS,
				AppNamespace: cfg.AppNS,
				Cores:        2,
			},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeveritySevere,
		},
		{
			Blueprint: node.MemoryHog{
				Namespace:        cfg.ChaosNS,
				AppNamespace:     cfg.AppNS,
				MemoryPercentage: 90,
			},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeveritySevere,
		},
		{
			Blueprint: node.IOStress{
				Namespace:             cfg.ChaosNS,
				AppNamespace:          cfg.AppNS,
				UtilizationPercentage: 90,
			},
			Scale:    metadata.ScaleNode,
			Severity: metadata.SeveritySevere,
		},
	}

	// For each scale (pod, part of deployment, entire deployment)
	for _, pods := range podsAffected {
		fs = append(fs, failures.Failure{
			Blueprint: container.NetworkLatency{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				NetworkLatency:         300,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeverityLight,
		}, failures.Failure{
			Blueprint: container.NetworkLatency{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				NetworkLatency:         3000,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeveritySevere,
		}, failures.Failure{
			Blueprint: container.CPUHog{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				Cores:                  1,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeveritySevere,
		}, failures.Failure{
			Blueprint: container.MemoryHog{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				MemoryConsumption:      1000,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeveritySevere,
		}, failures.Failure{
			Blueprint: pod.IOStress{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				UtilizationPercentage:  90,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeveritySevere,
		}, failures.Failure{
			Blueprint: pod.Delete{
				Namespace:              cfg.ChaosNS,
				AppNamespace:           cfg.AppNS,
				Interval:               1,
				Force:                  false,
				PodsAffectedPercentage: pods.Percentage,
			},
			Scale:    pods.Scale,
			Severity: metadata.SeverityCritical,
		})

		// For each severity (10%, 90%)
		for _, activities := range activitiesAffected {
			fs = append(fs, failures.Failure{
				Blueprint: container.NetworkLoss{
					Namespace:              cfg.ChaosNS,
					AppNamespace:           cfg.AppNS,
					LossPercentage:         activities.Percentage,
					PodsAffectedPercentage: pods.Percentage,
				},
				Scale:    pods.Scale,
				Severity: activities.Severity,
			}, failures.Failure{
				Blueprint: container.NetworkCorruption{
					Namespace:              cfg.ChaosNS,
					AppNamespace:           cfg.AppNS,
					CorruptionPercentage:   activities.Percentage,
					PodsAffectedPercentage: pods.Percentage,
				},
				Scale:    pods.Scale,
				Severity: activities.Severity,
			}, failures.Failure{
				Blueprint: container.NetworkDuplication{
					Namespace:              cfg.ChaosNS,
					AppNamespace:           cfg.AppNS,
					DuplicationPercentage:  activities.Percentage,
					PodsAffectedPercentage: pods.Percentage,
				},
				Scale:    pods.Scale,
				Severity: activities.Severity,
			})
		}
	}

	return fs
}
