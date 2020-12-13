package generator

import "time"

type Params struct {
	Stages        int
	Seed          int64
	StageDuration time.Duration
}

type Generator interface {
	Generate(params Params) (Scenario, error)
}
