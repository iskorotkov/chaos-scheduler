package generator

import "github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"

type Failure struct {
	Preset   experiments.Preset
	Scale    Scale
	Severity Severity
}
