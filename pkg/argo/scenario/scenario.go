package scenario

import (
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/input"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"strings"
	"text/template"
)

var (
	StagesError                = errors.New("can't create scenario with stages <= 0")
	TemplatesError             = errors.New("couldn't read templates")
	InsufficientTemplatesError = errors.New("can't create scenario without templates")
	TemplateParseError         = errors.New("couldn't parse text template")
	TemplateExecuteError       = errors.New("couldn't execute template text")
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
		return nil, StagesError
	}

	templates, err := input.Load(c.Path)
	if err != nil {
		logger.Error(err)
		return nil, TemplatesError
	}

	if len(templates) == 0 {
		return nil, InsufficientTemplatesError
	}

	stages := make([]Stage, 0)
	for i := 0; i < c.Stages; i++ {
		stage, err := createStage(i, templates)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	return Stage{step}, nil
}

func createStep(y string, ctx Context) (Step, error) {
	t, err := template.New(ctx.Name).Parse(y)
	if err != nil {
		logger.Error(err)
		return Step{}, TemplateParseError
	}

	b := &strings.Builder{}
	err = t.Execute(b, ctx)
	if err != nil {
		logger.Error(err)
		return Step{}, TemplateExecuteError
	}

	return Step{ctx.Name, b.String()}, nil
}
