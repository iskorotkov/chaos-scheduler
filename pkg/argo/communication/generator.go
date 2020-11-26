package communication

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/scenario"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/templates"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FormatConfig struct {
	TemplatePath string
	Scenario     scenario.Scenario
}

func GenerateWorkflow(config FormatConfig) (string, error) {
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

	spec["templates"], err = createTemplatesList(config.Scenario)
	if err != nil {
		return "", fmt.Errorf("couldn't create list of templates: %v", err)
	}

	res, err := yaml.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal workflow to yaml: %v", err)
	}

	return string(res), nil
}

func createTemplatesList(s scenario.Scenario) ([]interface{}, error) {
	actions := []interface{}{templates.NewStepsTemplate(s)}

	for _, stage := range s {
		for _, a := range stage {
			template, err := templates.NewManifestTemplate(a.Filename, a.Yaml, templates.ActionCreate)
			if err != nil {
				return nil, fmt.Errorf("couldn't create template: %v", err)
			}

			actions = append(actions, template)
		}
	}

	return actions, nil
}
