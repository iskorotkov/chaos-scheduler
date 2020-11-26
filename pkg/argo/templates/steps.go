package templates

import (
	"fmt"
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

	for i, stage := range scenario {
		newStage := make([]Step, 0)

		for j, action := range stage {
			friendlyName := fmt.Sprintf("%d.%d %s", i+1, j+1, action.Name)
			newStage = append(newStage, Step{Name: friendlyName, Template: action.Name})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
