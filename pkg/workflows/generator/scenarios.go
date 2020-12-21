package generator

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"time"
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
