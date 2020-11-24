package scenario

import (
	"fmt"
	"math/rand"
)

type Action struct {
	name       string
	definition string
}

type Stage []Action

type Scenario []Stage

type Failure interface {
	Name() string
	Yaml() string
}

type Config struct {
	Failures []Failure
	Stages   int
	Rng      *rand.Rand
}

func NewScenario(c Config) (Scenario, error) {
	if c.Stages <= 0 {
		return nil, fmt.Errorf("can't create scenario with stages <= 0")
	}

	if len(c.Failures) == 0 {
		return nil, fmt.Errorf("can't create scenario without experiments")
	}

	if c.Rng == nil {
		c.Rng = rand.New(rand.NewSource(42))
	}

	stages := make([]Stage, 0)
	for i := 0; i < c.Stages; i++ {
		stages = append(stages, createStage(c, i))
	}

	return stages, nil
}

func createStage(c Config, i int) Stage {
	f := c.Failures[(i % len(c.Failures))]
	return []Action{{f.Name(), f.Yaml()}}
}
