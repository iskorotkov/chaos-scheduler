package templates

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/assemblers"
)

type Step struct {
	Name     string `yaml:"name" json:"name"`
	Template string `yaml:"template" json:"template"`
}

type StepsTemplate struct {
	Name  string   `yaml:"name" json:"name"`
	Steps [][]Step `yaml:"steps" json:"steps"`
}

func NewStepsTemplate(scenario assemblers.Scenario) StepsTemplate {
	res := StepsTemplate{"entry", make([][]Step, 0)}

	for _, stage := range scenario {
		newStage := make([]Step, 0)

		for _, action := range stage {
			newStage = append(newStage, Step{action.Id(), action.Template()})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
