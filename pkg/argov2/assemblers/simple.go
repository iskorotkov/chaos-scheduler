package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"io/ioutil"
	"strings"
	"text/template"
)

type SimpleAssembler struct {
	WorkflowTemplate string
}

func (s SimpleAssembler) Assemble(scenario scenarios.Scenario) (Workflow, error) {
	if len(scenario.Stages()) == 0 {
		return nil, StagesError
	}

	workflowTemplate, err := ioutil.ReadFile(s.WorkflowTemplate)
	if err != nil {
		logger.Error(err)
		return nil, WorkflowTemplateError
	}

	workflow, err := marshall.FromYaml(workflowTemplate)
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

func createTemplatesList(s scenarios.Scenario) ([]interface{}, error) {
	actions := []interface{}{templates.NewStepsTemplate(s, generateActionId)}

	for i, stage := range s.Stages() {
		if len(stage.Actions()) == 0 {
			return nil, ActionsError
		}

		for j, action := range stage.Actions() {
			executedTemplate, err := executeTemplate(action.Template(), context{
				Name:     action.Name(),
				Duration: action.Duration(),
				Stage:    i,
				Index:    j,
			})
			if err != nil {
				return nil, err
			}

			id := generateActionId(action, i, j)
			manifestTemplate := templates.NewManifestTemplate(id, executedTemplate)

			actions = append(actions, manifestTemplate)
		}
	}

	return actions, nil
}

func generateActionId(action scenarios.Action, stage int, index int) string {
	return fmt.Sprintf("%s-%d-%d", action.Name(), stage+1, index+1)
}

func executeTemplate(content string, ctx context) (string, error) {
	t, err := template.New(ctx.Name).Parse(content)
	if err != nil {
		logger.Error(err)
		return "", TemplateParseError
	}

	b := &strings.Builder{}
	err = t.Execute(b, ctx)
	if err != nil {
		logger.Error(err)
		return "", TemplateExecuteError
	}

	return b.String(), nil
}
