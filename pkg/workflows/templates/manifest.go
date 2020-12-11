package templates

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

func NewManifestTemplate(name string, manifest string) Template {
	return Template{
		Name: name,
		Resource: &v1alpha1.ResourceTemplate{
			Action:   "apply",
			Manifest: manifest,
		},
	}
}
