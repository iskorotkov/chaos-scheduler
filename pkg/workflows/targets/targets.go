// Package targets manages fetching and managing list of targets.
package targets

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

var (
	ErrClient        = errors.New("couldn't create client")
	ErrFetch         = errors.New("couldn't fetch list of targets")
	ErrInvalidTarget = errors.New("target is not valid")
)

// Target describes potential target.
type Target struct {
	// Pod is a pod name.
	Pod string `json:"pod"`
	// Node is a node name where the pod runs.
	Node string `json:"node"`
	// MainContainer is a container to kill in failures.
	MainContainer string `json:"mainContainer"`
	// Containers is a list of all containers in the pod.
	Containers []string `json:"containers"`
	// AppLabel is a pod label.
	AppLabel string `json:"appLabel"`
	// AppLabelValue is a value of AppLabel (i.e. without a key part).
	AppLabelValue string
	// Labels is a list of labels from target metadata.
	Labels map[string]string `json:"labels"`
	// Annotations is a list of annotations from target metadata.
	Annotations map[string]string `json:"annotations"`
}

func (t Target) Generate(r *rand.Rand, _ int) reflect.Value {
	randomStr := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, r.Int())
	}

	containers := make([]string, 0)
	for i := 0; i < 1+r.Intn(10); i++ {
		containers = append(containers, randomStr("container"))
	}

	return reflect.ValueOf(Target{
		Pod:           randomStr("pod"),
		Node:          randomStr("node"),
		MainContainer: containers[r.Intn(len(containers))],
		Containers:    containers,
		AppLabel:      randomStr("label"),
		AppLabelValue: randomStr("label-value"),
		Labels: map[string]string{
			randomStr("label1"): randomStr("value"),
			randomStr("label2"): randomStr("value"),
		},
		Annotations: map[string]string{
			randomStr("annotation1"): randomStr("value"),
			randomStr("annotation2"): randomStr("value"),
		},
	})
}

// TargetFinder fetches targets.
type TargetFinder interface {
	// List returns fetched targets from specified namespace.
	List(namespace string, label string) ([]Target, error)
}

// TestTargetFinder mocks TargetFinder.
type TestTargetFinder struct {
	// Targets is a list of targets to return from List.
	Targets []Target
	// Err is an error to return from List.
	Err error
	// SubmittedNamespace is a namespace passed to List.
	SubmittedNamespace string
	// SubmittedLabel is a label passed to List.
	SubmittedLabel string
}

func (t TestTargetFinder) Generate(rand *rand.Rand, size int) reflect.Value {
	switch rand.Intn(10) {
	case 0:
		return reflect.ValueOf(TestTargetFinder{
			Targets: nil,
			Err:     ErrClient,
		})
	case 1:
		return reflect.ValueOf(TestTargetFinder{
			Targets: nil,
			Err:     ErrFetch,
		})
	default:
		var targets []Target
		for i := 0; i < rand.Intn(10); i++ {
			targets = append(targets, Target{}.Generate(rand, size).Interface().(Target))
		}

		return reflect.ValueOf(TestTargetFinder{
			Targets: targets,
			Err:     nil,
		})
	}
}

// List returns fake list of targets.
func (t *TestTargetFinder) List(namespace string, label string) ([]Target, error) {
	t.SubmittedNamespace = namespace
	t.SubmittedLabel = label
	return t.Targets, t.Err
}
