package scenarios

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/annotations"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroActions            = errors.New("can't create scenario out of 0 actions")
)

type TemplatedAction struct {
	Name        string
	Annotations annotations.List
	Template    string
}

type PlannedAction struct {
	Name     string
	Template string
}

type Stage struct {
	Actions  []PlannedAction
	Duration time.Duration
}

type Scenario struct {
	Stages []Stage
}
