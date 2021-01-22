package watcher

import (
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"time"
)

type Step struct {
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Phase       string            `json:"phase,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	StartedAt   time.Time         `json:"startedAt,omitempty"`
	FinishedAt  time.Time         `json:"finishedAt,omitempty"`
}

type Stage struct {
	Phase      string    `json:"phase,omitempty"`
	StartedAt  time.Time `json:"startedAt,omitempty"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Steps      []Step    `json:"steps,omitempty"`
}

type Event struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Type        string            `json:"type,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Phase       string            `json:"phase,omitempty"`
	StartedAt   time.Time         `json:"startedAt,omitempty"`
	FinishedAt  time.Time         `json:"finishedAt,omitempty"`
	Stages      []Stage           `json:"stages,omitempty"`
}

func buildNodesTree(ts []v1alpha1.Template, nodes v1alpha1.Nodes) ([]Stage, error) {
	stagesIDs, stepsIDs := splitStagesAndSteps(nodes)

	stages := make([]Stage, 0)
	for i := 0; i < len(stagesIDs); i++ {
		id := fmt.Sprintf("[%d]", i)
		stageStatus := stagesIDs[id]
		steps := make([]Step, 0)
		for _, stepID := range stageStatus.Children {
			stepStatus := stepsIDs[stepID]
			stepSpec, err := findStepSpec(ts, stepStatus)
			if err != nil {
				return nil, err
			}

			steps = append(steps, newStep(stepSpec.Metadata, stepStatus))
		}

		stages = append(stages, newStage(stageStatus, steps))
	}

	return stages, nil
}

func splitStagesAndSteps(nodes v1alpha1.Nodes) (map[string]v1alpha1.NodeStatus, map[string]v1alpha1.NodeStatus) {
	stagesID := make(map[string]v1alpha1.NodeStatus)
	stepsID := make(map[string]v1alpha1.NodeStatus)

	for _, n := range nodes {
		if n.Type == "StepGroup" {
			stagesID[n.DisplayName] = n
		} else {
			stepsID[n.ID] = n
		}
	}
	return stagesID, stepsID
}

func findStepSpec(ts []v1alpha1.Template, stepStatus v1alpha1.NodeStatus) (templates.Template, error) {
	for _, t := range ts {
		if t.Name == stepStatus.TemplateName {
			return templates.Template(t), nil
		}
	}

	return templates.Template{}, SpecError
}

func newStep(metadata v1alpha1.Metadata, n v1alpha1.NodeStatus) Step {
	return Step{
		Name:        n.TemplateName,
		Type:        string(n.Type),
		Phase:       string(n.Phase),
		Labels:      metadata.Labels,
		Annotations: metadata.Annotations,
		StartedAt:   n.StartedAt.Time,
		FinishedAt:  n.FinishedAt.Time,
	}
}

func newStage(n v1alpha1.NodeStatus, steps []Step) Stage {
	return Stage{
		Phase:      string(n.Phase),
		StartedAt:  n.StartedAt.Time,
		FinishedAt: n.FinishedAt.Time,
		Steps:      steps,
	}
}

func newEvent(e *workflow.WorkflowWatchEvent) (*Event, error) {
	stages, err := buildNodesTree(e.Object.Spec.Templates, e.Object.Status.Nodes)
	if err != nil {
		return nil, err
	}

	return &Event{
		Name:        e.Object.Name,
		Namespace:   e.Object.Namespace,
		Type:        e.Type,
		Labels:      e.Object.Labels,
		Annotations: e.Object.Annotations,
		Phase:       string(e.Object.Status.Phase),
		StartedAt:   e.Object.Status.StartedAt.Time,
		FinishedAt:  e.Object.Status.FinishedAt.Time,
		Stages:      stages,
	}, nil
}
