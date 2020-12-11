package generators

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroActions            = errors.New("can't create scenario out of 0 actions")
)

type Action struct {
	Type   presets.ExperimentType
	Info   presets.Info
	Target targets.Target
	Engine presets.Engine
}

type Stage struct {
	Actions  []Action
	Duration time.Duration
}

type Scenario struct {
	Stages []Stage
}
