package assemble

import (
	"errors"
	"fmt"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/metadata"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"reflect"
	"time"
)

var (
	StagesError         = errors.New("number of stages must be positive")
	ActionsError        = errors.New("number of actions in every stage must be positive")
	ActionMarshallError = errors.New("couldn't marshall action to yaml")
	MetadataError       = errors.New("couldn't set metadata")
)

type ActionExtension interface {
	Apply(action generate.Action, stageIndex, actionIndex int) []templates.Template
}

type StageExtension interface {
	Apply(stage generate.Stage, stageIndex int) []templates.Template
}

type WorkflowExtension interface {
	Apply(ids [][]string) []templates.Template
}

type Extensions struct {
	Action   []ActionExtension
	Stage    []StageExtension
	Workflow []WorkflowExtension
}

func (e Extensions) Generate(r *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(Extensions{
		Action: []ActionExtension{
			// No action extensions implemented
		},
		Stage: []StageExtension{
			UseSuspend(),
			UseStageMonitor("stage-monitor", "target-ns", time.Duration(r.Intn(60)), &zap.SugaredLogger{}),
		},
		Workflow: []WorkflowExtension{
			UseSteps(),
		},
	})
}

func Assemble(scenario generate.Scenario, extensions Extensions) (templates.Workflow, error) {
	if len(scenario.Stages) == 0 {
		return templates.Workflow{}, StagesError
	}

	ts, err := createTemplatesList(scenario, extensions)
	if err != nil {
		return templates.Workflow{}, err
	}

	wf := templates.NewWorkflow("entry", ts,
		templates.WithNamespace("litmus"),
		templates.WithNamePrefix("workflow-"),
		templates.WithServiceAccount("argo-chaos"))

	return wf, nil
}

func createTemplatesList(scenario generate.Scenario, extensions Extensions) ([]templates.Template, error) {
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

			if err := addFailureMetadata(&manifestTemplate, action); err != nil {
				return nil, err
			}

			actions = append(actions, manifestTemplate)
			stageIDs = append(stageIDs, id)

			extensionsActions, extensionsIDs, err := applyActionExtensions(action, stageIndex, actionIndex, extensions.Action)
			if err != nil {
				return nil, err
			}

			actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)
		}

		extensionsActions, extensionsIDs, err := applyStageExtensions(stage, stageIndex, extensions.Stage)
		if err != nil {
			return nil, err
		}

		actions, stageIDs = append(actions, extensionsActions...), append(stageIDs, extensionsIDs...)

		ids = append(ids, stageIDs)
	}

	workflowActions, err := applyWorkflowExtensions(ids, extensions.Workflow)
	if err != nil {
		return nil, err
	}

	actions = append(actions, workflowActions...)
	return actions, nil
}

func addFailureMetadata(t *templates.Template, action generate.Action) error {
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

func addUtilityMetadata(t *templates.Template, severity api.Severity, scale api.Scale) error {
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

func applyWorkflowExtensions(ids [][]string, extensions []WorkflowExtension) ([]templates.Template, error) {
	actions := make([]templates.Template, 0)

	if extensions != nil {
		for _, extension := range extensions {
			createdExtensions := extension.Apply(ids)
			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)
			}
		}
	}

	for i := 0; i < len(actions); i++ {
		if err := addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func applyStageExtensions(stage generate.Stage, stageIndex int, extensions []StageExtension) ([]templates.Template, []string, error) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if extensions != nil {
		for _, extension := range extensions {
			createdExtensions := extension.Apply(stage, stageIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.Id())
				}
			}
		}
	}

	for i := 0; i < len(actions); i++ {
		if err := addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, nil, err
		}
	}

	return actions, stageIDs, nil
}

func applyActionExtensions(action generate.Action, stageIndex, actionIndex int, extensions []ActionExtension) ([]templates.Template, []string, error) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if extensions != nil {
		for _, extension := range extensions {
			createdExtensions := extension.Apply(action, stageIndex, actionIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.Id())
				}
			}
		}
	}

	for i := 0; i < len(actions); i++ {
		if err := addUtilityMetadata(&actions[i], api.SeverityHarmless, api.ScaleCluster); err != nil {
			return nil, nil, err
		}
	}

	return actions, stageIDs, nil
}
