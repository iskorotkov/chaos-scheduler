package assemble

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type steps struct{}

func (s steps) Apply(ids [][]string) []templates.Template {
	return []templates.Template{
		templates.NewStepsTemplate(ids),
	}
}

func UseSteps() WorkflowExtension {
	return steps{}
}
