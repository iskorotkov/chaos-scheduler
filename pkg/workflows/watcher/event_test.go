package watcher

import (
	"errors"
	"fmt"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

type nodes v1alpha1.Nodes

func (n nodes) Generate(r *rand.Rand, _ int) reflect.Value {
	rs := func(s string) string {
		return fmt.Sprintf("%s-%d", s, r.Int())
	}

	nodeTypes := []v1alpha1.NodeType{v1alpha1.NodeTypeDAG, v1alpha1.NodeTypeSteps, v1alpha1.NodeTypeSuspend}
	nodePhases := []v1alpha1.NodePhase{v1alpha1.NodeSucceeded, v1alpha1.NodeFailed, v1alpha1.NodePending, v1alpha1.NodeRunning}

	statuses := make(nodes)
	for i := 0; i < r.Intn(10); i++ {
		statuses[rs("name")] = v1alpha1.NodeStatus{
			ID:            rs("id"),
			Name:          rs("name"),
			DisplayName:   rs("display-name"),
			Type:          nodeTypes[r.Intn(len(nodeTypes))],
			TemplateName:  rs("template-name"),
			TemplateScope: rs("template-scope"),
			Phase:         nodePhases[r.Intn(len(nodePhases))],
			BoundaryID:    rs("boundary-id"),
			Message:       rs("message"),
			StartedAt:     v1.Time{Time: time.Now().Add(-5 * time.Hour)},
			FinishedAt:    v1.Time{Time: time.Now().Add(-20 * time.Minute)},
			PodIP:         rs("pod-ip"),
		}
	}

	return reflect.ValueOf(statuses)
}

func validateTime(start time.Time, finish time.Time) error {
	if !start.Before(finish) {
		return errors.New("start time must be before finish time")
	}

	if finish.Sub(start) > 2*time.Hour {
		return errors.New("workflow must be executed for 2 hours at most")
	}

	return nil
}

func Test_buildNodesTree(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	f := func(ts []v1alpha1.Template, n nodes) bool {
		stages, err := buildNodesTree(ts, v1alpha1.Nodes(n))
		if err != nil {
			t.Log(err)
			return false
		}

		for _, stage := range stages {
			if stage.Phase == "" {
				t.Log("stage phase must not be empty")
				return false
			}

			if err := validateTime(stage.StartedAt, stage.FinishedAt); err != nil {
				t.Log(err)
				return false
			}

			for _, step := range stage.Steps {
				if err := validateTime(step.StartedAt, step.FinishedAt); err != nil {
					t.Log(err)
					return false
				}

				if step.Name == "" ||
					step.Phase == "" ||
					step.Type == "" {
					t.Log("step name, phase and type must not be empty")
					return false
				}
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
