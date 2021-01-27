package generate

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
)

type Cost float64

type Budget struct {
	MaxFailures int
	MaxPoints   Cost
}

func DefaultBudget() Budget {
	return Budget{MaxFailures: 3, MaxPoints: 12}
}

type Modifiers struct {
	ByScale    map[metadata.Scale]Cost
	BySeverity map[metadata.Severity]Cost
}

func DefaultModifiers() Modifiers {
	return Modifiers{
		ByScale: map[metadata.Scale]Cost{
			metadata.ScaleContainer:      1,
			metadata.ScalePod:            1,
			metadata.ScaleDeploymentPart: 1.5,
			metadata.ScaleDeployment:     2,
			metadata.ScaleNode:           4,
		},
		BySeverity: map[metadata.Severity]Cost{
			metadata.SeverityLight:    1,
			metadata.SeveritySevere:   1.5,
			metadata.SeverityCritical: 2,
		},
	}
}

func calculateCost(modifiers Modifiers, f failures.Failure) Cost {
	severity, ok := modifiers.BySeverity[f.Severity]
	if !ok {
		severity = 1
	}

	scale, ok := modifiers.ByScale[f.Scale]
	if !ok {
		scale = 1
	}

	return severity * scale
}
