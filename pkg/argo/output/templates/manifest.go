package templates

import "fmt"

const (
	ActionCreate = Action("create")
	ActionDelete = Action("delete")
)

type Action string

type Resource struct {
	Action   Action `yaml:"action"`
	Manifest string `yaml:"manifest"`
}

type ManifestTemplate struct {
	Name     string   `yaml:"name"`
	Resource Resource `yaml:"resource"`
}

func NewManifestTemplate(name string, manifest string, action Action) (ManifestTemplate, error) {
	if name == "" {
		return ManifestTemplate{}, fmt.Errorf("can't create template with empty name")
	}

	if action != ActionCreate && action != ActionDelete {
		return ManifestTemplate{}, fmt.Errorf("only 'create' and 'delete' actions are supported")
	}

	if manifest == "" {
		return ManifestTemplate{}, fmt.Errorf("can't create template with empty manifest")
	}

	return ManifestTemplate{name, Resource{action, manifest}}, nil
}
