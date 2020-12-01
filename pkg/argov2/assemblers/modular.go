package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"io/ioutil"
	"strings"
	"text/template"
)

type ModularAssembler struct {
	WorkflowTemplate   string
	ActionExtensions   []extensions.ActionExtension
	StageExtensions    []extensions.StageExtension
	WorkflowExtensions []extensions.WorkflowExtension
}

func (a ModularAssembler) Assemble(scenario scenarios.Scenario) (Workflow, error) {
	if len(scenario.Stages()) == 0 {
		return nil, StagesError
	}

	workflowTemplate, err := ioutil.ReadFile(a.WorkflowTemplate)
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

	spec["templates"], err = a.createTemplatesList(scenario)
	if err != nil {
		return nil, err
	}

	return Workflow(workflow), nil
}

func DefaultModularAssembler(workflowTemplate string) Assembler {
	return ModularAssembler{
		WorkflowTemplate:   workflowTemplate,
		ActionExtensions:   []extensions.ActionExtension{},
		StageExtensions:    []extensions.StageExtension{extensions.Suspend{}},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.Steps{}},
	}
}

func (a ModularAssembler) createTemplatesList(s scenarios.Scenario) ([]interface{}, error) {
	actions := make([]interface{}, 0)
	ids := make([][]string, 0)

	for stageIndex, stage := range s.Stages() {
		if len(stage.Actions()) == 0 {
			return nil, ActionsError
		}

		stageIds := make([]string, 0)

		for actionIndex, action := range stage.Actions() {
			executedTemplate, err := executeTemplate(action.Template(), context{
				Name:     action.Name(),
				Duration: stage.Duration(),
				Stage:    stageIndex,
				Index:    actionIndex,
			})
			if err != nil {
				return nil, err
			}

			id := fmt.Sprintf("%s-%d-%d", action.Name(), stageIndex+1, actionIndex+1)
			manifestTemplate := templates.NewManifestTemplate(id, executedTemplate)

			actions = append(actions, manifestTemplate)
			stageIds = append(stageIds, id)

			// Apply action extensions
			if a.ActionExtensions != nil {
				for _, ext := range a.ActionExtensions {
					createdExt := ext.Apply(action, stageIndex, actionIndex)
					if createdExt != nil {
						actions = append(actions, createdExt)
						stageIds = append(stageIds, createdExt.Id())
					}
				}
			}
		}

		// Apply stage extensions
		if a.StageExtensions != nil {
			for _, ext := range a.StageExtensions {
				createdExt := ext.Apply(stage, stageIndex)
				if createdExt != nil {
					actions = append(actions, createdExt)
					stageIds = append(stageIds, createdExt.Id())
				}
			}
		}

		ids = append(ids, stageIds)
	}

	// Apply workflow extensions
	if a.WorkflowExtensions != nil {
		for _, ext := range a.WorkflowExtensions {
			createdExt := ext.Apply(ids)
			if createdExt != nil {
				actions = append(actions, createdExt)
			}
		}
	}

	return actions, nil
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
