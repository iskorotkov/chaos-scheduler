package generator

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
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
	Name     string            `json:"name"`
	Severity metadata.Severity `json:"severity"`
	Scale    metadata.Scale    `json:"scale"`
	Engine   failures.Engine   `json:"engine"`
	Target   targets.Target    `json:"target"`
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
