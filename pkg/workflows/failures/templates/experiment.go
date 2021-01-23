package templates

import (
	"fmt"
	"math/rand"
	"reflect"
)

type ExperimentName string

type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type ExperimentComponents struct {
	Env []EnvVar `json:"env,omitempty" yaml:"env,omitempty"`
}

type ExperimentSpec struct {
	Components ExperimentComponents `json:"components" yaml:"components"`
}

type Experiment struct {
	Name string         `json:"name" yaml:"name"`
	Spec ExperimentSpec `json:"spec" yaml:"spec"`
}

func (e Experiment) Generate(r *rand.Rand, _ int) reflect.Value {
	rs := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, r.Int())
	}

	return reflect.ValueOf(NewExperiment(ExperimentParams{
		Name: rs("name"),
		Env: map[string]string{
			rs("env1"): rs("value"),
			rs("env2"): rs("value"),
		},
	}))
}

type ExperimentParams struct {
	Name string
	Env  map[string]string
}

func NewExperiment(params ExperimentParams) Experiment {
	envVarList := make([]EnvVar, 0)
	for k, v := range params.Env {
		envVarList = append(envVarList, EnvVar{Name: k, Value: v})
	}

	return Experiment{
		Name: params.Name,
		Spec: ExperimentSpec{
			Components: ExperimentComponents{
				Env: envVarList,
			},
		},
	}
}
