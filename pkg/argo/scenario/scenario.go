package scenario

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/input"
	"strings"
	"text/template"
)

type Context struct {
	Stage int
	Step  int
	Name  string
}

type Step struct {
	Name string
	Yaml string
}

type Stage []Step

type Scenario []Stage

type Config struct {
	Path   string
	Stages int
	Seed   int64
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

	stages := make([]Stage, 0)
	for i := 0; i < c.Stages; i++ {
		stage, err := createStage(i, templates)
		if err != nil {
			return nil, fmt.Errorf("couldn't create stage: %v", err)
		}

		stages = append(stages, stage)
	}

	return stages, nil
}

func createStage(stage int, t []input.Template) (Stage, error) {
	selected := t[(stage % len(t))]

	name := fmt.Sprintf("%s-%s-%v-%v", "cluster", selected.Filename, stage+1, 1)
	ctx := Context{Stage: stage, Name: name}

	step, err := createStep(selected.Yaml, ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't update template yaml: %v", err)
	}

	return Stage{step}, nil
}

func createStep(y string, ctx Context) (Step, error) {
	t, err := template.New(ctx.Name).Parse(y)
	if err != nil {
		return Step{}, fmt.Errorf("couldn't parse text template: %v", err)
	}

	b := &strings.Builder{}
	err = t.Execute(b, ctx)
	if err != nil {
		return Step{}, fmt.Errorf("couldn't execute template text: %v", err)
	}

	return Step{ctx.Name, b.String()}, nil
}
