package templates

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

type Step struct {
	Name     string `yaml:"name" json:"name"`
	Template string `yaml:"template" json:"template"`
}

type StepsTemplate struct {
	Name  string   `yaml:"name" json:"name"`
	Steps [][]Step `yaml:"steps" json:"steps"`
}

type IdGenerator func(action scenarios.Action, stage int, index int) string

func NewStepsTemplate(scenario scenarios.Scenario, generator IdGenerator) StepsTemplate {
	res := StepsTemplate{"entry", make([][]Step, 0)}

	for i, stage := range scenario.Stages() {
		newStage := make([]Step, 0)

		for j, action := range stage.Actions() {
			id := generator(action, i, j)
			newStage = append(newStage, Step{id, id})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
