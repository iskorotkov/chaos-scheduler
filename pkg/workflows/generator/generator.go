package generator

import (
	"errors"
	"time"
)

var (
	TargetsError = errors.New("couldn't get list of targets")
)

type Params struct {
	Stages        int
	Seed          int64
	StageDuration time.Duration
}

type Generator interface {
	Generate(params Params) (Scenario, error)
}
