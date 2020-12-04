package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
)

type ScenarioParams struct {
	Stages int
	Seed   int64
}

type Generator interface {
	Generate(presetsList presets.List, params ScenarioParams) (Scenario, error)
}
