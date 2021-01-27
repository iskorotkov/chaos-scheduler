package templates

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type Template v1alpha1.Template

func (t Template) Id() string {
	return t.Name
}
