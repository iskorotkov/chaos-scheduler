// Package assemble converts generate.Scenario to Workflow.
package assemble

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/rx"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"github.com/iskorotkov/metadata"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrStages             = errors.New("scenario must contain at least one stage")
	ErrActions            = errors.New("stage must contain at least one action")
	ErrActionMarshalError = errors.New("couldn't marshal action to yaml")
	ErrMetadata           = errors.New("couldn't set template metadata")
)

// Workflow is the definition of a workflow resource.
type Workflow v1alpha1.Workflow

func (w Workflow) Generate(rand *rand.Rand, size int) reflect.Value {
	rs := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, rand.Intn(100))
	}

	var ts []v1alpha1.Template
	for i := 0; i <= rand.Intn(10); i++ {
		randomTemplate := templates.Template{}.Generate(rand, size).Interface().(templates.Template)
		ts = append(ts, v1alpha1.Template(randomTemplate))
	}

	return reflect.ValueOf(Workflow{
		TypeMeta: v1.TypeMeta{
			Kind:       rs("kind"),
			APIVersion: rs("api-version"),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:         rs("name"),
			GenerateName: rs("generate-name"),
			Namespace:    rs("namespace"),
			Labels:       rx.Rmap(rand, 10),
			Annotations:  rx.Rmap(rand, 10),
		},
		Spec: v1alpha1.WorkflowSpec{
			Templates:          ts,
			Entrypoint:         rs("entrypoint"),
			ServiceAccountName: rs("sa-name"),
		},
	})
}

type Option func(wf *Workflow)

func WithNamespace(ns string) Option {
	return func(wf *Workflow) {
		wf.ObjectMeta.Namespace = ns
	}
}

func WithNamePrefix(prefix string) Option {
	return func(wf *Workflow) {
		wf.ObjectMeta.GenerateName = prefix
	}
}

func WithServiceAccount(sa string) Option {
	return func(wf *Workflow) {
		wf.Spec.ServiceAccountName = sa
	}
}

func NewWorkflow(entrypoint string, templates []templates.Template, opts ...Option) Workflow {
	argoTemplates := make([]v1alpha1.Template, 0)
	for _, template := range templates {
		argoTemplates = append(argoTemplates, v1alpha1.Template(template))
	}

	wf := Workflow{
		TypeMeta: v1.TypeMeta{
			Kind:       "Workflow",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint: entrypoint,
			Templates:  argoTemplates,
		},
	}

	for _, opt := range opts {
		opt(&wf)
	}

	return wf
}

// ActionExt appends templates to a stage based on a previous action.
type ActionExt interface {
	Apply(action generate.Action, stageIndex, actionIndex int) []templates.Template
}

// StageExt appends templates to a stage based on other actions of the stage.
type StageExt interface {
	Apply(stage generate.Stage, stageIndex int) []templates.Template
}

// WorkflowExt appends templates to a workflow based on actions IDs.
type WorkflowExt interface {
	Apply(ids [][]string) []templates.Template
}

type ExtCollection struct {
	Action   []ActionExt
	Stage    []StageExt
	Workflow []WorkflowExt
}

func (e ExtCollection) Generate(r *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(ExtCollection{
		Action: []ActionExt{
			// No action extensions implemented
		},
		Stage: []StageExt{
			UseMonitor("stage-monitor", "target-ns", time.Duration(r.Intn(60)), &zap.SugaredLogger{}),
		},
		Workflow: []WorkflowExt{
			UseSteps(),
		},
	})
}

// Assemble converts generate.Scenario to Workflow.
func Assemble(scenario generate.Scenario, extensions ExtCollection) (Workflow, error) {
	if len(scenario.Stages) == 0 {
		return Workflow{}, ErrStages
	}

	ts, err := createTemplatesList(scenario, extensions)
	if err != nil {
		return Workflow{}, err
	}

	wf := NewWorkflow("entry", ts,
		WithNamespace("litmus"),
		WithNamePrefix("workflow-"),
		WithServiceAccount("argo-chaos"))

	return wf, nil
}

// createTemplatesList returns linear list of all templates.
func createTemplatesList(scenario generate.Scenario, extensions ExtCollection) ([]templates.Template, error) {
	actions := make([]templates.Template, 0)
	ids := make([][]string, 0)

	for stageIndex, stage := range scenario.Stages {
		if len(stage.Actions) == 0 {
			return nil, ErrActions
		}

		stageIDs := make([]string, 0)

		for actionIndex, action := range stage.Actions {
			id := fmt.Sprintf("%s-%d-%d", action.Name, stageIndex+1, actionIndex+1)

			action.Engine.Metadata.Name = id
			manifest, err := yaml.Marshal(action.Engine)
			if err != nil {
				return nil, ErrActionMarshalError
			}

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

// addFailureMetadata adds metadata to a failure template.
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
		return ErrMetadata
	}

	t.Metadata.Labels, t.Metadata.Annotations = objectMeta.Labels, objectMeta.Annotations
	return nil
}

// addUtilityMetadata adds metadata to a utility template.
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
		return ErrMetadata
	}

	t.Metadata.Labels, t.Metadata.Annotations = objectMeta.Labels, objectMeta.Annotations
	return nil
}

func applyWorkflowExtensions(ids [][]string, extensions []WorkflowExt) ([]templates.Template, error) {
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

func applyStageExtensions(stage generate.Stage, stageIndex int, extensions []StageExt) ([]templates.Template, []string, error) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if extensions != nil {
		for _, extension := range extensions {
			createdExtensions := extension.Apply(stage, stageIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.ID())
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

func applyActionExtensions(action generate.Action, stageIndex, actionIndex int, extensions []ActionExt) ([]templates.Template, []string, error) {
	actions := make([]templates.Template, 0)
	stageIDs := make([]string, 0)

	if extensions != nil {
		for _, extension := range extensions {
			createdExtensions := extension.Apply(action, stageIndex, actionIndex)

			if createdExtensions != nil {
				actions = append(actions, createdExtensions...)

				for _, created := range createdExtensions {
					stageIDs = append(stageIDs, created.ID())
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
