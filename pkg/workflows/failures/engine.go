package failures

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppInfo struct {
	AppNS    string `json:"appns" yaml:"appns"`
	AppLabel string `json:"applabel" yaml:"applabel"`
	AppKind  string `json:"appkind" yaml:"appkind"`
}

type EngineSpec struct {
	AppInfo             AppInfo      `json:"appinfo" yaml:"appinfo"`
	JobCleanupPolicy    string       `json:"jobCleanupPolicy,omitempty" yaml:"jobCleanupPolicy,omitempty"`
	Monitoring          bool         `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	AnnotationsCheck    bool         `json:"annotationsCheck,omitempty" yaml:"annotationsCheck,omitempty"`
	EngineState         string       `json:"engineState,omitempty" yaml:"engineState,omitempty"`
	ChaosServiceAccount string       `json:"chaosServiceAccount,omitempty" yaml:"chaosServiceAccount,omitempty"`
	Experiments         []Experiment `json:"experiments" yaml:"experiments"`
}

type Engine struct {
	APIVersion string        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string        `json:"kind" yaml:"kind"`
	Metadata   v1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       EngineSpec    `json:"spec" yaml:"spec"`
}

type EngineParams struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
	AppInfo     AppInfo
	Experiments []Experiment
}

func NewEngine(params EngineParams) Engine {
	return Engine{
		Kind:       "ChaosEngine",
		APIVersion: "litmuschaos.io/v1alpha1",
		Metadata: v1.ObjectMeta{
			Name:        params.Name,
			Namespace:   params.Namespace,
			Labels:      params.Labels,
			Annotations: params.Annotations,
		},
		Spec: EngineSpec{
			AppInfo:             params.AppInfo,
			JobCleanupPolicy:    "delete",
			Monitoring:          false,
			AnnotationsCheck:    true,
			EngineState:         "active",
			ChaosServiceAccount: "litmus-admin",
			Experiments:         params.Experiments,
		},
	}
}
