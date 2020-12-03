package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/scenarios"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"io/ioutil"
)

type ModularAssembler struct {
	WorkflowTemplate   string
	ActionExtensions   []extensions.ActionExtension
	StageExtensions    []extensions.StageExtension
	WorkflowExtensions []extensions.WorkflowExtension
}

func (a ModularAssembler) Assemble(scenario scenarios.Scenario) (Workflow, error) {
	if len(scenario.Stages) == 0 {
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

func NewModularAssembler(template string, ae []extensions.ActionExtension, se []extensions.StageExtension, we []extensions.WorkflowExtension) Assembler {
	return ModularAssembler{
		WorkflowTemplate:   template,
		ActionExtensions:   ae,
		StageExtensions:    se,
		WorkflowExtensions: we,
	}
}

func (a ModularAssembler) createTemplatesList(scenario scenarios.Scenario) ([]interface{}, error) {
	actions := make([]interface{}, 0)
	ids := make([][]string, 0)

	for stageIndex, stage := range scenario.Stages {
		if len(stage.Actions) == 0 {
			return nil, ActionsError
		}

		stageIds := make([]string, 0)

		for actionIndex, action := range stage.Actions {
			manifest, err := marshall.ToYaml(action.Engine)
			if err != nil {
				return nil, ActionMarshallError
			}

			id := fmt.Sprintf("%s-%d-%d", action.Type, stageIndex+1, actionIndex+1)
			manifestTemplate := templates.NewManifestTemplate(id, string(manifest))

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
