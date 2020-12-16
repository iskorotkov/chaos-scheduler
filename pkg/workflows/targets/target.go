package targets

import "fmt"

type Target struct {
	Pod           string
	Deployment    string
	Node          string
	Containers    []string
	Labels        map[string]string
	Annotations   map[string]string
	SelectorLabel string
}

func (t Target) MainContainer() string {
	return t.Containers[0]
}

func (t Target) Selector() string {
	return fmt.Sprintf("%s=%s", t.SelectorLabel, t.Labels[t.SelectorLabel])
}
