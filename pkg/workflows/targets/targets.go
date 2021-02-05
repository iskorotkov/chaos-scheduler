package targets

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

var (
	ConfigError = errors.New("couldn't read config")
	ClientError = errors.New("couldn't create client from config")
	FetchError  = errors.New("couldn't fetch info from Kubernetes")
)

type Target struct {
	Pod           string            `json:"pod"`
	Deployment    string            `json:"deployment"`
	Node          string            `json:"node"`
	MainContainer string            `json:"mainContainer"`
	Containers    []string          `json:"containers"`
	AppLabel      string            `json:"appLabel"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
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
		Deployment:    randomStr("deploy"),
		Node:          randomStr("node"),
		MainContainer: containers[r.Intn(len(containers))],
		Containers:    containers,
		AppLabel:      randomStr("label"),
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

type TargetFinder interface {
	List(namespace string, label string) ([]Target, error)
}

type TestTargetFinder struct {
	Targets            []Target
	Err                error
	SubmittedNamespace string
	SubmittedLabel     string
}

func (t TestTargetFinder) Generate(rand *rand.Rand, size int) reflect.Value {
	switch rand.Intn(10) {
	case 0:
		return reflect.ValueOf(TestTargetFinder{
			Targets: nil,
			Err:     ClientError,
		})
	case 1:
		return reflect.ValueOf(TestTargetFinder{
			Targets: nil,
			Err:     ConfigError,
		})
	case 2:
		return reflect.ValueOf(TestTargetFinder{
			Targets: nil,
			Err:     FetchError,
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

func (t *TestTargetFinder) List(namespace string, label string) ([]Target, error) {
	t.SubmittedNamespace = namespace
	t.SubmittedLabel = label
	return t.Targets, t.Err
}
