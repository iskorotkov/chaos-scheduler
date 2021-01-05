package templates

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Template v1alpha1.Template

func (t Template) Id() string {
	return t.Name
}

type Workflow v1alpha1.Workflow

type Option func(wf *Workflow)

func NewWorkflow(namespace, generateName, entrypoint, serviceAccountName string, templates []Template, opts ...Option) Workflow {
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
			Namespace:    namespace,
			GenerateName: generateName,
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint:         entrypoint,
			ServiceAccountName: serviceAccountName,
			Templates:          argoTemplates,
		},
	}

	for _, opt := range opts {
		opt(&wf)
	}

	return wf
}

func WithLabel(key, value string) Option {
	return func(wf *Workflow) {
		wf.ObjectMeta.Labels[key] = value
	}
}

func WithAnnotation(key, value string) Option {
	return func(wf *Workflow) {
		wf.ObjectMeta.Annotations[key] = value
	}
}
