package assemblers

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"gopkg.in/yaml.v2"
)

type ModularAssembler struct {
	Extensions extensions.List
}

func (a ModularAssembler) Assemble(scenario generator.Scenario) (templates.Workflow, error) {
	if len(scenario.Stages) == 0 {
		return templates.Workflow{}, StagesError
	}

	ts, err := a.createTemplatesList(scenario)
	if err != nil {
		return templates.Workflow{}, err
	}

	wf := templates.NewWorkflow("litmus", "workflow-", "entry", "argo-chaos", ts)

	return wf, nil
}

func NewModularAssembler(ext extensions.List) Assembler {
	return ModularAssembler{Extensions: ext}
}

func (a ModularAssembler) createTemplatesList(scenario generator.Scenario) ([]templates.Template, error) {
	actions := make([]templates.Template, 0)
	ids := make([][]string, 0)

	for stageIndex, stage := range scenario.Stages {
		if len(stage.Actions) == 0 {
			return nil, ActionsError
		}

		stageIds := make([]string, 0)

		for actionIndex, action := range stage.Actions {
			manifest, err := yaml.Marshal(action.Engine)
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
					createdExtensions := ext.Apply(action, stageIndex, actionIndex)

					if createdExtensions != nil {
						actions = append(actions, createdExtensions...)

						for _, created := range createdExtensions {
							stageIds = append(stageIds, created.Id())
						}
					}
				}
			}
		}

		// Apply stage extensions
		if a.Extensions.StageExtensions != nil {
			for _, ext := range a.Extensions.StageExtensions {
				createdExtensions := ext.Apply(stage, stageIndex)

				if createdExtensions != nil {
					actions = append(actions, createdExtensions...)

					for _, created := range createdExtensions {
						stageIds = append(stageIds, created.Id())
					}
				}
			}
		}

		ids = append(ids, stageIds)
	}

	// Apply workflow extensions
	if a.Extensions.WorkflowExtensions != nil {
		for _, ext := range a.Extensions.WorkflowExtensions {
			createdExtensions := ext.Apply(ids)
			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)
			}
		}
	}

	return actions, nil
}
