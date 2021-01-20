package assemblers

import (
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/metadata"
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

	wf := templates.NewWorkflow("entry", ts,
		templates.WithNamespace("litmus"),
		templates.WithNamePrefix("workflow-"),
		templates.WithServiceAccount("argo-chaos"))

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

		stageIDs := make([]string, 0)

		for actionIndex, action := range stage.Actions {
			meta := api.TemplateMetadata{
				Version:  api.VersionV1,
				Type:     api.TypeFailure,
				Severity: action.Info.Severity,
				Scale:    action.Info.Scale,
			}

			if err := metadata.Marshal(&action.Engine.Metadata, &meta, api.Prefix); err != nil {
				return nil, ActionMarshallError
			}

			manifest, err := yaml.Marshal(action.Engine)
			if err != nil {
				return nil, ActionMarshallError
			}

			id := fmt.Sprintf("%s-%d-%d", action.Info.Name, stageIndex+1, actionIndex+1)
			manifestTemplate := templates.NewManifestTemplate(id, string(manifest))

			actions = append(actions, manifestTemplate)
			stageIDs = append(stageIDs, id)

			extensionsActions, extensionsIDs := a.applyActionExtensions(action, stageIndex, actionIndex)
			actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)
		}

		extensionsActions, extensionsIDs := a.applyStageExtensions(stage, stageIndex)
		actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)

		ids = append(ids, stageIDs)
	}

	actions = append(actions, a.applyWorkflowExtensions(ids)...)
	return actions, nil
}

func (a ModularAssembler) applyWorkflowExtensions(ids [][]string) []templates.Template {
	actions := make([]templates.Template, 0)

	if a.Extensions.WorkflowExtensions != nil {
		for _, ext := range a.Extensions.WorkflowExtensions {
			createdExtensions := ext.Apply(ids)
			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)
			}
		}
	}
	return actions
}

func (a ModularAssembler) applyStageExtensions(stage generator.Stage, stageIndex int) ([]templates.Template, []string) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if a.Extensions.StageExtensions != nil {
		for _, ext := range a.Extensions.StageExtensions {
			createdExtensions := ext.Apply(stage, stageIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.Id())
				}
			}
		}
	}
	return actions, stageIDs
}

func (a ModularAssembler) applyActionExtensions(action generator.Action, stageIndex int, actionIndex int) ([]templates.Template, []string) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if a.Extensions.ActionExtensions != nil {
		for _, ext := range a.Extensions.ActionExtensions {
			createdExtensions := ext.Apply(action, stageIndex, actionIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.Id())
				}
			}
		}
	}
	return actions, stageIDs
}
