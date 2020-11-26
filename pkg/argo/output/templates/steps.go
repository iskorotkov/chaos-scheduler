package templates

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/scenario"
)

type Step struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
}

type StepsTemplate struct {
	Name  string   `yaml:"name"`
	Steps [][]Step `yaml:"steps"`
}

func NewStepsTemplate(scenario scenario.Scenario) StepsTemplate {
	res := StepsTemplate{"entry", make([][]Step, 0)}

	for _, stage := range scenario {
		newStage := make([]Step, 0)

		for _, template := range stage {
			newStage = append(newStage, Step{template.StepName, template.TemplateName})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
