package advanced

import (
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
)

type Cost float64

type Budget struct {
	MaxFailures int
	MaxPoints   Cost
}

type Modifiers struct {
	ByScale    map[metadata.Scale]Cost
	BySeverity map[metadata.Severity]Cost
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
