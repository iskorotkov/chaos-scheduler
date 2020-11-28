package templates

import (
	"errors"
)

var (
	FilenameError = errors.New("can't create template with empty name")
	ManifestError = errors.New("can't create template with empty manifest")
)

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
		return ManifestTemplate{}, FilenameError
	}

	if manifest == "" {
		return ManifestTemplate{}, ManifestError
	}

	return ManifestTemplate{name, Resource{"create", manifest}}, nil
}
