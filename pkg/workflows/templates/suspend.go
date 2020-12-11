package templates

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"time"
)

func NewSuspendTemplate(name string, duration time.Duration) Template {
	return Template{
		Name: name,
		Suspend: &v1alpha1.SuspendTemplate{
			Duration: duration.String(),
		},
	}
}
