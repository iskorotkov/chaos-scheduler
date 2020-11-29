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
}

type action struct {
	name     string
	template string
}

func (a action) Name() string {
	return a.name
}

func (a action) Template() string {
	return a.template
}

type Stage interface {
	Actions() []Action
	Duration() time.Duration
}

type stage struct {
	actions  []Action
	duration time.Duration
}

func (s stage) Actions() []Action {
	return s.actions
}

func (s stage) Duration() time.Duration {
	return s.duration
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
