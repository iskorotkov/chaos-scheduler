package scenario

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/input"
	"math/rand"
)

type Template struct {
	StepName     string
	TemplateName string
	Yaml         string
}

type Stage []Template

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

	templates, err := input.Load(c.Path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read templates: %v", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("can't create scenario without templates")
	}

	if c.Rng == nil {
		c.Rng = rand.New(rand.NewSource(42))
	}

	stages := make([]Stage, 0)
	for i := 0; i < c.Stages; i++ {
		stages = append(stages, createStage(i, templates))
	}

	return stages, nil
}

func createStage(stage int, t []input.Template) Stage {
	selected := t[(stage % len(t))]

	name := fmt.Sprintf("%s-%s-%v-%v", "cluster", selected.Filename, stage+1, 1)
	template := Template{name, name, selected.Yaml}

	return Stage{template}
}
