package output

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/scenario"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"io/ioutil"
)

var (
	TemplatePathError      = errors.New("template path wasn't provided")
	ScenarioError          = errors.New("scenario wasn't provided")
	TemplateReadError      = errors.New("couldn't read template file")
	TemplateMarshalError   = errors.New("couldn't marshall template from yaml")
	TemplateUnmarshalError = errors.New("couldn't unmarshall template from yaml")
	TemplatePropertyError  = errors.New("couldn't find required template property")
	ManifestTemplateError  = errors.New("couldn't create manifest template")
)

type Config struct {
	TemplatePath string
	Scenario     scenario.Scenario
}

func GenerateFromConfig(config Config) (string, error) {
	if config.TemplatePath == "" {
		return "", TemplatePathError
	}

	if config.Scenario == nil {
		return "", ScenarioError
	}

	template, err := ioutil.ReadFile(config.TemplatePath)
	if err != nil {
		logger.Error(err)
		return "", TemplateReadError
	}

	workflow, err := marshall.FromYaml(template)
	if err != nil {
		logger.Error(err)
		return "", TemplateUnmarshalError
	}

	spec, ok := workflow["spec"].(marshall.Tree)
	if !ok {
		logger.Warning("couldn't access property 'spec'")
		return "", TemplatePropertyError
	}

	spec["templates"], err = createTemplatesList(config.Scenario)
	if err != nil {
		return "", err
	}

	msg := struct {
		Workflow marshall.Tree `json:"workflow"`
	}{workflow}

	res, err := marshall.ToJson(msg)
	if err != nil {
		logger.Error(err)
		return "", TemplateMarshalError
	}

	return string(res), nil
}

func createTemplatesList(s scenario.Scenario) ([]interface{}, error) {
	actions := []interface{}{templates.NewStepsTemplate(s)}

	for _, stage := range s {
		for _, a := range stage {
			template, err := templates.NewManifestTemplate(a.Name, a.Yaml)
			if err != nil {
				logger.Error(err)
				return nil, ManifestTemplateError
			}

			actions = append(actions, template)
		}
	}

	return actions, nil
}
