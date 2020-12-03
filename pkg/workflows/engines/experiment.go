package engines

type ExperimentType string

type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type ExperimentComponents struct {
	Env []EnvVar `json:"env" yaml:"env"`
}

type ExperimentSpec struct {
	Components ExperimentComponents `json:"components" yaml:"components"`
}

type Experiment struct {
	Name ExperimentType `json:"name" yaml:"name"`
	Spec ExperimentSpec `json:"spec" yaml:"spec"`
}

type ExperimentParams struct {
	Type ExperimentType
	Env  map[string]string
}

func NewExperiment(params ExperimentParams) Experiment {
	envVarList := make([]EnvVar, 0)
	for k, v := range params.Env {
		envVarList = append(envVarList, EnvVar{Name: k, Value: v})
	}

	return Experiment{
		Name: params.Type,
		Spec: ExperimentSpec{
			Components: ExperimentComponents{
				Env: envVarList,
			},
		},
	}
}
