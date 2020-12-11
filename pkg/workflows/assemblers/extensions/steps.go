package extensions

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type Steps struct{}

func (s Steps) Apply(ids [][]string) []templates.Template {
	return []templates.Template{
		templates.NewStepsTemplate(ids),
	}
}

func UseSteps() WorkflowExtension {
	return Steps{}
}
