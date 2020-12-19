package generator

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
)

type Action struct {
	Info   experiments.Info
	Target targets.Target
	Engine experiments.Engine
}

type Stage struct {
	Actions  []Action
	Duration time.Duration
}

type Scenario struct {
	Stages []Stage
}
