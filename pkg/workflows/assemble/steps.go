package assemble

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
)

type steps struct{}

// Apply adds steps template to the workflow
func (s steps) Apply(ids [][]string) []templates.Template {
	return []templates.Template{
		templates.NewStepsTemplate(ids),
	}
}

// UseSteps returns a workflow extension that adds steps template to the workflow.
func UseSteps() WorkflowExt {
	return steps{}
}
