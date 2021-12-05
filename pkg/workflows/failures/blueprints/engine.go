package blueprints

import (
	"fmt"
	"math/rand"
	"reflect"
)

// AppInfo describes failure target.
type AppInfo struct {
	// AppNS describes target namespace.
	AppNS string `json:"appns" yaml:"appns"`
	// AppLabel describes target label.
	AppLabel string `json:"applabel" yaml:"applabel"`
	// AppKind describes target kind.
	AppKind string `json:"appkind" yaml:"appkind"`
}

type EngineMetadata struct {
	Name        string            `json:"name,omitempty" yaml:"name"`
	Namespace   string            `json:"namespace,omitempty" yaml:"namespace"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations"`
}

// EngineSpec describes engine spec.
type EngineSpec struct {
	// AppInfo describes failure target.
	AppInfo AppInfo `json:"appinfo" yaml:"appinfo"`
	// JobCleanUpPolicy describes cleanup actions after the engine finishes execution.
	JobCleanUpPolicy string `json:"jobCleanUpPolicy,omitempty" yaml:"jobCleanUpPolicy,omitempty"`
	// Monitoring describes whether monitoring data should be exported.
	Monitoring bool `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	// AnnotationCheck describes whether failures must influence only pods with the annotation set.
	AnnotationCheck string `json:"annotationCheck,omitempty" yaml:"annotationCheck,omitempty"`
	// EngineState is a current state of the engine.
	EngineState string `json:"engineState,omitempty" yaml:"engineState,omitempty"`
	// ChaosServiceAccount is a ServiceAccount name to use to cause chaos.
	ChaosServiceAccount string `json:"chaosServiceAccount,omitempty" yaml:"chaosServiceAccount,omitempty"`
	// Experiments is a list of experiments included in the engine.
	Experiments []Experiment `json:"experiments" yaml:"experiments"`
}

// Engine is a set of experiments and associated values.
type Engine struct {
	APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   EngineMetadata `json:"metadata" yaml:"metadata"`
	Spec       EngineSpec     `json:"spec" yaml:"spec"`
}

type EngineParams struct {
	Name string
	// Namespace is a namespace to create Engine in.
	Namespace string
	// Labels is a list of Engine labels.
	Labels map[string]string
	// Annotations is a list of Engine annotations.
	Annotations map[string]string
	// AppInfo is a target of the Engine.
	AppInfo AppInfo
	// Experiments is a list of chaos failures to cause.
	Experiments []Experiment
}

func (e Engine) Generate(r *rand.Rand, size int) reflect.Value {
	rs := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, r.Int())
	}

	experiments := make([]Experiment, 0)
	for i := 0; i < 1+r.Intn(10); i++ {
		experiments = append(experiments, Experiment{}.Generate(r, size).Interface().(Experiment))
	}

	return reflect.ValueOf(NewEngine(EngineParams{
		Name:      rs("name"),
		Namespace: rs("namespace"),
		Labels: map[string]string{
			rs("label1"): rs("value"),
			rs("label2"): rs("value"),
		},
		Annotations: map[string]string{
			rs("annotation1"): rs("value"),
			rs("annotation2"): rs("value"),
		},
		AppInfo: AppInfo{
			AppNS:    rs("app-ns"),
			AppLabel: rs("app-label"),
			AppKind:  rs("app-kind"),
		},
		Experiments: experiments,
	}))
}

func NewEngine(params EngineParams) Engine {
	return Engine{
		Kind:       "ChaosEngine",
		APIVersion: "litmuschaos.io/v1alpha1",
		Metadata: EngineMetadata{
			Name:        params.Name,
			Namespace:   params.Namespace,
			Labels:      params.Labels,
			Annotations: params.Annotations,
		},
		Spec: EngineSpec{
			AppInfo:             params.AppInfo,
			JobCleanUpPolicy:    "delete",
			Monitoring:          false,
			AnnotationCheck:     "true",
			EngineState:         "active",
			ChaosServiceAccount: "litmus-admin",
			Experiments:         params.Experiments,
		},
	}
}
