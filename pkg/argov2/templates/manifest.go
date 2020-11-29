package templates

type Resource struct {
	Action   string `yaml:"action" json:"action"`
	Manifest string `yaml:"manifest" json:"manifest"`
}

type ManifestTemplate struct {
	Name     string   `yaml:"name" json:"name"`
	Resource Resource `yaml:"resource" json:"resource"`
}

func NewManifestTemplate(name string, manifest string) ManifestTemplate {
	return ManifestTemplate{name, Resource{"apply", manifest}}
}
