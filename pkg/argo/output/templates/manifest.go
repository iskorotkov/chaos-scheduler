package templates

import "fmt"

type Resource struct {
	Action   string `yaml:"action"`
	Manifest string `yaml:"manifest"`
}

type ManifestTemplate struct {
	Name     string   `yaml:"name"`
	Resource Resource `yaml:"resource"`
}

func NewManifestTemplate(name string, manifest string) (ManifestTemplate, error) {
	if name == "" {
		return ManifestTemplate{}, fmt.Errorf("can't create template with empty name")
	}

	if manifest == "" {
		return ManifestTemplate{}, fmt.Errorf("can't create template with empty manifest")
	}

	return ManifestTemplate{name, Resource{"create", manifest}}, nil
}
