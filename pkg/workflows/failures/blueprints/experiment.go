package blueprints

import (
	"fmt"
	"math/rand"
	"reflect"
)

// EnvVar describes environment variable.
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

// Experiment is a smallest part of the failure that determines its effects.
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
	// Experiment is a special value used to determine the experiment type.
	Name string
	// Env is a list of env vars.
	Env map[string]string
}

// NewExperiment returns new Experiment.
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
