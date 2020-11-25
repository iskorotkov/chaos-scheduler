package argo

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FormatConfig struct {
	TemplatePath string
	Scenario     Scenario
}

func Format(config FormatConfig) (string, error) {
	if config.TemplatePath == "" {
		return "", fmt.Errorf("template path wasn't provided")
	}

	if config.Scenario == nil {
		return "", fmt.Errorf("scenario wasn't provided")
	}

	template, err := ioutil.ReadFile(config.TemplatePath)
	if err != nil {
		return "", fmt.Errorf("couldn't read template file: %v", err)
	}

	workflow := make(map[interface{}]interface{})
	err = yaml.Unmarshal(template, workflow)
	if err != nil {
		return "", fmt.Errorf("couldn't unmarshall template from yaml: %v", err)
	}

	spec, ok := workflow["spec"].(map[interface{}]interface{})
	if !ok {
		return "", fmt.Errorf("couldn't get 'spec' property of template")
	}

	spec["templates"] = createTemplatesList(config.Scenario)

	res, err := yaml.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal workflow to yaml: %v", err)
	}

	return string(res), nil
}

func createTemplatesList(s Scenario) []interface{} {
	entryAction := createEntryAction(s)
	actions := []interface{}{entryAction}

	for _, stage := range s {
		for _, a := range stage {
			actions = append(actions, a.Yaml)
		}
	}

	return actions
}
