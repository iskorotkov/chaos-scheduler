package targets

import (
	"fmt"
	"math/rand"
	"reflect"
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
