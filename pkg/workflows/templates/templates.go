// Package templates allows to use different types of templates in generated workflows.
package templates

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/rx"
	v1 "k8s.io/api/core/v1"
)

// Template is a reusable and composable unit of execution in a workflow.
type Template v1alpha1.Template

func (t Template) Generate(rand *rand.Rand, _ int) reflect.Value {
	rs := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, rand.Intn(100))
	}

	return reflect.ValueOf(Template{
		Name: rs("name"),
		Metadata: v1alpha1.Metadata{
			Annotations: rx.Rmap(rand, 10),
			Labels:      rx.Rmap(rand, 10),
		},
		Steps: []v1alpha1.ParallelSteps{
			{[]v1alpha1.WorkflowStep{
				{
					Name:     rx.Rstr(rand, "name"),
					Template: rx.Rstr(rand, "template"),
				}, {
					Name:     rx.Rstr(rand, "name"),
					Template: rx.Rstr(rand, "template"),
				}, {
					Name:     rx.Rstr(rand, "name"),
					Template: rx.Rstr(rand, "template"),
				},
			}}, {[]v1alpha1.WorkflowStep{
				{
					Name:     rx.Rstr(rand, "name"),
					Template: rx.Rstr(rand, "template"),
				}, {
					Name:     rx.Rstr(rand, "name"),
					Template: rx.Rstr(rand, "template"),
				},
			}},
		},
		Container: &v1.Container{
			Name:  rx.Rstr(rand, "name"),
			Image: rx.Rstr(rand, "image"),
		},
		ServiceAccountName: rs("sa-name"),
	})
}

// ID is a template unique ID.
func (t Template) ID() string {
	return t.Name
}
