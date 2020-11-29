package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"io/ioutil"
)

type SimpleAssembler struct {
	WorkflowTemplate string
}

func (s SimpleAssembler) Assemble(scenario Scenario) (Workflow, error) {
	if len(scenario) == 0 {
		return nil, StagesError
	}

	template, err := ioutil.ReadFile(s.WorkflowTemplate)
	if err != nil {
		logger.Error(err)
		return nil, WorkflowTemplateError
	}

	workflow, err := marshall.FromYaml(template)
	if err != nil {
		logger.Error(err)
		return nil, WorkflowTemplateUnmarshalError
	}

	spec, ok := workflow["spec"].(marshall.Tree)
	if !ok {
		logger.Warning("couldn't access property 'spec'")
		return nil, WorkflowTemplatePropertyError
	}

	spec["templates"], err = createTemplatesList(scenario)
	if err != nil {
		return nil, err
	}

	return Workflow(workflow), nil
}

func NewSimpleAssembler(workflowTemplate string) Assembler {
	return SimpleAssembler{workflowTemplate}
}

func createTemplatesList(s Scenario) ([]interface{}, error) {
	actions := []interface{}{templates.NewStepsTemplate(s)}

	for i, stage := range s {
		if len(stage) == 0 {
			return nil, ActionsError
		}

		for j, action := range stage {

			id := fmt.Sprintf("%s-%d-%d", action.Id(), i+1, j+1)
			template := templates.NewManifestTemplate(id, action.Template())

			actions = append(actions, template)
		}
	}

	return actions, nil
}
