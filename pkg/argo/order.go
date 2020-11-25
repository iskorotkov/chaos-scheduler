package argo

import "fmt"

type actionReference struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
}

type entryAction struct {
	Name  string              `yaml:"name"`
	Steps [][]actionReference `yaml:"steps"`
}

func createEntryAction(scenario Scenario) entryAction {
	res := entryAction{"entry", make([][]actionReference, 0)}

	for i, stage := range scenario {
		newStage := make([]actionReference, 0)

		for j, action := range stage {
			friendlyName := fmt.Sprintf("%d.%d %s", i+1, j+1, action.Name)
			newStage = append(newStage, actionReference{Name: friendlyName, Template: action.Name})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
