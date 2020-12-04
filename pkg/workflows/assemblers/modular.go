package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/marshall"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/workflow"
)

type ModularAssembler struct {
	Extensions extensions.List
}

func (a ModularAssembler) Assemble(scenario generators.Scenario) (workflow.Workflow, error) {
	if len(scenario.Stages) == 0 {
		return workflow.Workflow{}, StagesError
	}

	ts, err := a.createTemplatesList(scenario)
	if err != nil {
		return workflow.Workflow{}, err
	}

	wf := workflow.NewWorkflow("litmus", "workflow-", "entry", "argo-chaos", ts)

	return wf, nil
}

func NewModularAssembler(ext extensions.List) Assembler {
	return ModularAssembler{Extensions: ext}
}

func (a ModularAssembler) createTemplatesList(scenario generators.Scenario) ([]workflow.Template, error) {
	actions := make([]workflow.Template, 0)
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
			if a.Extensions.ActionExtensions != nil {
				for _, ext := range a.Extensions.ActionExtensions {
					createdExt := ext.Apply(action, stageIndex, actionIndex)
					if createdExt != nil {
						actions = append(actions, createdExt)
						stageIds = append(stageIds, createdExt.Id())
					}
				}
			}
		}

		// Apply stage extensions
		if a.Extensions.StageExtensions != nil {
			for _, ext := range a.Extensions.StageExtensions {
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
	if a.Extensions.WorkflowExtensions != nil {
		for _, ext := range a.Extensions.WorkflowExtensions {
			createdExt := ext.Apply(ids)
			if createdExt != nil {
				actions = append(actions, createdExt)
			}
		}
	}

	return actions, nil
}
