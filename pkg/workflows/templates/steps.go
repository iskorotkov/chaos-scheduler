package templates

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

// NewStepsTemplate returns a template representing a list of stages of steps, where each step contains one other template.
func NewStepsTemplate(ids [][]string) Template {
	parallel := make([]v1alpha1.ParallelSteps, 0)

	for _, stage := range ids {
		newStage := make([]v1alpha1.WorkflowStep, 0)

		for _, id := range stage {
			newStage = append(newStage, v1alpha1.WorkflowStep{Name: id, Template: id})
		}

		parallel = append(parallel, v1alpha1.ParallelSteps{Steps: newStage})
	}

	return Template{
		Name:  "entry",
		Steps: parallel,
	}
}
