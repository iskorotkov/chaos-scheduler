package scenarios

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/engines"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroActions            = errors.New("can't create scenario out of 0 actions")
)

type PlannedAction struct {
	Type   engines.ExperimentType
	Engine engines.Engine
}

type Stage struct {
	Actions  []PlannedAction
	Duration time.Duration
}

type Scenario struct {
	Stages []Stage
}
