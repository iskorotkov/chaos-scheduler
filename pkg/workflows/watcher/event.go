package watcher

import (
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"time"
)

type Step struct {
	Name       string    `json:"name,omitempty"`
	Type       string    `json:"type,omitempty"`
	Phase      string    `json:"phase,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

type Stage struct {
	Phase      string    `json:"phase,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	Steps      []Step    `json:"steps,omitempty"`
}

type Event struct {
	Name       string            `json:"name,omitempty"`
	Namespace  string            `json:"namespace,omitempty"`
	Type       string            `json:"type,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	Phase      string            `json:"phase,omitempty"`
	StartedAt  time.Time         `json:"started_at,omitempty"`
	FinishedAt time.Time         `json:"finished_at,omitempty"`
	Stages     []Stage           `json:"stages,omitempty"`
}

func buildNodesTree(nodes v1alpha1.Nodes) []Stage {
	stagesID := make(map[string]v1alpha1.NodeStatus)
	stepsID := make(map[string]v1alpha1.NodeStatus)

	for _, n := range nodes {
		if n.Type == "StepGroup" {
			stagesID[n.DisplayName] = n
		} else {
			stepsID[n.ID] = n
		}
	}

	stages := make([]Stage, 0)
	for i := 0; i < len(stagesID); i++ {
		stage := stagesID[fmt.Sprintf("[%d]", i)]
		steps := make([]Step, 0)
		for _, stepID := range stage.Children {
			step := stepsID[stepID]
			steps = append(steps, newStep(step))
		}

		stages = append(stages, newStage(stage, steps))
	}

	return stages
}

func newStep(n v1alpha1.NodeStatus) Step {
	return Step{
		Name:       n.TemplateName,
		Type:       string(n.Type),
		Phase:      string(n.Phase),
		StartedAt:  n.StartedAt.Time,
		FinishedAt: n.FinishedAt.Time,
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

func newEvent(e *workflow.WorkflowWatchEvent) *Event {
	return &Event{
		Name:       e.Object.Name,
		Namespace:  e.Object.Namespace,
		Type:       e.Type,
		Labels:     e.Object.Labels,
		Phase:      string(e.Object.Status.Phase),
		StartedAt:  e.Object.Status.StartedAt.Time,
		FinishedAt: e.Object.Status.FinishedAt.Time,
		Stages:     buildNodesTree(e.Object.Status.Nodes),
	}
}
