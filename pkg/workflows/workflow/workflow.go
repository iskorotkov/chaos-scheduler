package workflow

type Template interface{}

type Metadata struct {
	Namespace    string `json:"namespace" yaml:"namespace"`
	GenerateName string `json:"generateName,omitempty" yaml:"generateName,omitempty"`
}

type Spec struct {
	Entrypoint         string     `json:"entrypoint" yaml:"entrypoint"`
	ServiceAccountName string     `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
	Templates          []Template `json:"templates" yaml:"templates"`
}

type Workflow struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
}

func NewWorkflow(namespace string, generateName string, entrypoint string, serviceAccountName string, templates []Template) Workflow {
	return Workflow{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Workflow",
		Metadata: Metadata{
			Namespace:    namespace,
			GenerateName: generateName,
		},
		Spec: Spec{
			Entrypoint:         entrypoint,
			ServiceAccountName: serviceAccountName,
			Templates:          templates,
		},
	}
}
