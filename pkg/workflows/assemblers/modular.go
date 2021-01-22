package assemblers

import (
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/metadata"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			if err := a.addFailureMetadata(action); err != nil {
				return nil, err
			}

			manifest, err := yaml.Marshal(action.Engine)
			if err != nil {
				return nil, ActionMarshallError
			}

			id := fmt.Sprintf("%s-%d-%d", action.Name, stageIndex+1, actionIndex+1)
			manifestTemplate := templates.NewManifestTemplate(id, string(manifest))

			actions = append(actions, manifestTemplate)
			stageIDs = append(stageIDs, id)

			extensionsActions, extensionsIDs, err := a.applyActionExtensions(action, stageIndex, actionIndex)
			if err != nil {
				return nil, err
			}

			actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)
		}

		extensionsActions, extensionsIDs, err := a.applyStageExtensions(stage, stageIndex)
		if err != nil {
			return nil, err
		}

		actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)

		ids = append(ids, stageIDs)
	}

	workflowActions, err := a.applyWorkflowExtensions(ids)
	if err != nil {
		return nil, err
	}

	actions = append(actions, workflowActions...)
	return actions, nil
}

func (a ModularAssembler) addFailureMetadata(action generator.Action) error {
	values := api.TemplateMetadata{
		Version:  api.VersionV1,
		Type:     api.TypeFailure,
		Severity: action.Severity,
		Scale:    action.Scale,
	}

	if err := metadata.Marshal(&action.Engine.Metadata, &values, api.Prefix); err != nil {
		return MetadataError
	}
	return nil
}

func (a ModularAssembler) addUtilityMetadata(action templates.Template, severity api.Severity, scale api.Scale) error {
	values := api.TemplateMetadata{
		Version:  api.VersionV1,
		Type:     api.TypeUtility,
		Severity: severity,
		Scale:    scale,
	}

	var objectMeta v1.ObjectMeta
	if err := metadata.Marshal(&objectMeta, &values, api.Prefix); err != nil {
		return MetadataError
	}

	action.Metadata.Labels, action.Metadata.Labels = objectMeta.Labels, objectMeta.Annotations
	return nil
}

func (a ModularAssembler) applyWorkflowExtensions(ids [][]string) ([]templates.Template, error) {
	actions := make([]templates.Template, 0)

	if a.Extensions.WorkflowExtensions != nil {
		for _, ext := range a.Extensions.WorkflowExtensions {
			createdExtensions := ext.Apply(ids)
			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)
			}
		}
	}

	for _, action := range actions {
		if err := a.addUtilityMetadata(action, api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func (a ModularAssembler) applyStageExtensions(stage generator.Stage, stageIndex int) ([]templates.Template, []string, error) {
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

	for _, action := range actions {
		if err := a.addUtilityMetadata(action, api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, nil, err
		}
	}

	return actions, stageIDs, nil
}

func (a ModularAssembler) applyActionExtensions(action generator.Action, stageIndex, actionIndex int) ([]templates.Template, []string, error) {
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

	for _, action := range actions {
		if err := a.addUtilityMetadata(action, api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, nil, err
		}
	}

	return actions, stageIDs, nil
}
