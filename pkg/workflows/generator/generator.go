package generator

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroFailures           = errors.New("can't create scenario out of 0 failures")
	TargetsError           = errors.New("couldn't get list of targets")
)

type Action struct {
	Info   experiments.Info   `json:"info"`
	Target targets.Target     `json:"target"`
	Engine experiments.Engine `json:"engine"`
}

type Stage struct {
	Actions  []Action      `json:"actions"`
	Duration time.Duration `json:"duration"`
}

type Scenario struct {
	Stages []Stage `json:"stages"`
}

type Params struct {
	Stages        int
	Seed          int64
	StageDuration time.Duration
}

type Generator interface {
	Generate(params Params) (Scenario, error)
}
