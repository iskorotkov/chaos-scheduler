package extensions

import "github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"

type Steps struct{}

func (s Steps) Apply(ids [][]string) Extension {
	return templates.NewStepsTemplate(ids)
}

func UseSteps() WorkflowExtension {
	return Steps{}
}
