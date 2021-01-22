package assemblers

import (
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/metadata"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ModularAssembler struct {
	Extensions extensions.List
	logger     *zap.SugaredLogger
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

func NewModularAssembler(ext extensions.List, logger *zap.SugaredLogger) Assembler {
	return ModularAssembler{Extensions: ext, logger: logger}
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
			manifest, err := yaml.Marshal(action.Engine)
			if err != nil {
				return nil, ActionMarshallError
			}

			id := fmt.Sprintf("%s-%d-%d", action.Name, stageIndex+1, actionIndex+1)
			manifestTemplate := templates.NewManifestTemplate(id, string(manifest))

			if err := a.addFailureMetadata(&manifestTemplate, action); err != nil {
				return nil, err
			}

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

func (a ModularAssembler) addFailureMetadata(t *templates.Template, action generator.Action) error {
	values := api.TemplateMetadata{
		Version:  api.VersionV1,
		Type:     api.TypeFailure,
		Severity: action.Severity,
		Scale:    action.Scale,
	}

	// TODO: Do not use temporary ObjectMeta to marshal data to
	objectMeta := v1.ObjectMeta{
		Labels:      t.Metadata.Labels,
		Annotations: t.Metadata.Annotations,
	}
	if err := metadata.Marshal(&objectMeta, &values, api.Prefix); err != nil {
		return MetadataError
	}

	t.Metadata.Labels, t.Metadata.Annotations = objectMeta.Labels, objectMeta.Annotations
	return nil
}

func (a ModularAssembler) addUtilityMetadata(t *templates.Template, severity api.Severity, scale api.Scale) error {
	values := api.TemplateMetadata{
		Version:  api.VersionV1,
		Type:     api.TypeUtility,
		Severity: severity,
		Scale:    scale,
	}

	// TODO: Do not use temporary ObjectMeta to marshal data to
	objectMeta := v1.ObjectMeta{
		Labels:      t.Metadata.Labels,
		Annotations: t.Metadata.Annotations,
	}
	if err := metadata.Marshal(&objectMeta, &values, api.Prefix); err != nil {
		return MetadataError
	}

	t.Metadata.Labels, t.Metadata.Annotations = objectMeta.Labels, objectMeta.Annotations
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

	for i := 0; i < len(actions); i++ {
		if err := a.addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
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

	for i := 0; i < len(actions); i++ {
		if err := a.addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
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

	for i := 0; i < len(actions); i++ {
		if err := a.addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, nil, err
		}
	}

	return actions, stageIDs, nil
}
