package argo

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

type Config struct {
	Path   string
	Stages int
	Rng    *rand.Rand
}

func NewScenario(c Config) (Scenario, error) {
	if c.Stages <= 0 {
		return nil, fmt.Errorf("can't create scenario with stages <= 0")
	}

	failures, err := Load(c.Path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read failures: %v", err)
	}

	if len(failures) == 0 {
		return nil, fmt.Errorf("can't create scenario without experiments")
	}

	if c.Rng == nil {
		c.Rng = rand.New(rand.NewSource(42))
	}

	stages := make([]Stage, 0)
	for i := 0; i < c.Stages; i++ {
		stages = append(stages, createStage(failures, i))
	}

	return stages, nil
}

func createStage(failures []failure, i int) Stage {
	f := failures[(i % len(failures))]
	return []Action{{f.name, f.yaml}}
}
