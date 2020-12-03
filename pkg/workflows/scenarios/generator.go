package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"
)

type ScenarioParams struct {
	Stages int
	Seed   int64
}

type Generator interface {
	Generate(factories []engines.Factory, params ScenarioParams) (Scenario, error)
}
