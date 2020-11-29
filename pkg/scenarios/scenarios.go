package scenarios

import (
	"errors"
	"time"
)

var (
	NonPositiveStagesError = errors.New("number of stages must be positive")
	TooManyStagesError     = errors.New("number of stages can't be that high")
	ZeroActions            = errors.New("can't create scenario out of 0 actions")
)

type Template interface {
	Name() string
	Template() string
}

type Action interface {
	Name() string
	Template() string
	Duration() time.Duration
}

type action struct {
	name     string
	template string
	duration time.Duration
}

func (a action) Name() string {
	return a.name
}

func (a action) Template() string {
	return a.template
}

func (a action) Duration() time.Duration {
	return a.duration
}

type Stage interface {
	Actions() []Action
}

type stage []Action

func (s stage) Actions() []Action {
	return s
}

type Scenario interface {
	Stages() []Stage
}

type scenario []Stage

func (s scenario) Stages() []Stage {
	return s
}

type Config struct {
	Stages int
	Seed   int64
}

type Generator interface {
	Generate(templates []Template, config Config) (Scenario, error)
}
