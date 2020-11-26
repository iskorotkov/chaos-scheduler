package scenario

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/io"
	"math/rand"
)

type Stage []io.Template

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

	failures, err := io.Load(c.Path)
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

func createStage(actions []io.Template, stage int) Stage {
	template := actions[(stage % len(actions))]
	return Stage{template}
}
