// Package config loads app config.
package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"math/rand"
	"reflect"
	"time"
)

var (
	ErrParse = errors.New("couldn't parse config from env vars")
)

// Config describes app settings read from env vars.
type Config struct {
	// Infrastructure.

	// ArgoServer is an address of Argo used for executing generated workflows.
	ArgoServer string `env:"ARGO_SERVER"`
	// StageMonitorImage is a Docker image used for monitoring the state of the system under test.
	StageMonitorImage string `env:"STAGE_MONITOR_IMAGE"`

	// Target.

	// AppNS is a namespace of the system under test.
	AppNS string `env:"APP_NS"`
	// AppLabel is a metadata label used for target selection.
	AppLabel string `env:"APP_LABEL"`
	// ChaosNS is a namespace where workflows will be created.
	ChaosNS string `env:"CHAOS_NS"`
	// Development describes whether is in development or not.
	Development bool `env:"DEVELOPMENT"`

	// Workflow

	// StageDuration describes a duration of each test stage.
	StageDuration time.Duration `env:"STAGE_DURATION"`
	// StageInterval describes an interval between test stages.
	StageInterval time.Duration `env:"STAGE_INTERVAL"`

	// Node.

	// NodeCPUHogCores is a number of node cores to occupy.
	NodeCPUHogCores         int `env:"NODE_CPU_HOG_CORES"`
	// NodeMemoryHogPercentage is a percent of total node RAM to occupy.
	NodeMemoryHogPercentage int `env:"NODE_MEMORY_HOG_PERCENTAGE"`
	// NodeIOStressPercentage is a percent of total IO bandwidth of node to occupy.
	NodeIOStressPercentage  int `env:"NODE_IO_STRESS_PERCENTAGE"`

	// Pod.

	// PodIOStressPercentage is a percent of total IO bandwidth of a pod to occupy.
	PodIOStressPercentage int `env:"POD_IO_STRESS_PERCENTAGE"`

	// Container.

	// ContainerCPUHogCores is a number of container cores to occupy
	ContainerCPUHogCores int `env:"CONTAINER_CPU_HOG_CORES"`
	// ContainerMemoryHogMB is a percent of total container RAM to occupy.
	ContainerMemoryHogMB int `env:"CONTAINER_MEMORY_HOG_MB"`

	// Deployment part.

	// DeploymentPartPodsPercentage is a percent of all deployment pods to be affected.
	DeploymentPartPodsPercentage int `env:"DEPLOYMENT_PART_PODS_PERCENTAGE"`

	// Severity.

	// LightSeverityPercentage is a percentage of pods affected.
	LightSeverityPercentage  int `env:"LIGHT_SEVERITY_PERCENTAGE"`
	// SevereSeverityPercentage is a percentage of pods affected.
	SevereSeverityPercentage int `env:"SEVERE_SEVERITY_PERCENTAGE"`

	// Latency.

	// LightNetworkLatencyMS is a network latency in ms to apply.
	LightNetworkLatencyMS  int `env:"LIGHT_NETWORK_LATENCY_MS"`
	// SevereNetworkLatencyMS is a network latency in ms to apply.
	SevereNetworkLatencyMS int `env:"SEVERE_NETWORK_LATENCY_MS"`

	// Pod delete.

	// PodDeleteInterval is an interval between two successive pod failures.
	PodDeleteInterval int  `env:"POD_DELETE_INTERVAL"`
	// PodDeleteForce indicates whether to use immediate forceful deletion (with 0s grace period).
	PodDeleteForce    bool `env:"POD_DELETE_FORCE"`
}

func (c Config) Generate(r *rand.Rand, _ int) reflect.Value {
	rs := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, r.Intn(100))
	}
	return reflect.ValueOf(Config{
		ArgoServer:        rs("argo-server"),
		StageMonitorImage: rs("stage-monitor-image"),
		AppNS:             rs("app-ns"),
		AppLabel:          rs("app-label"),
		ChaosNS:           rs("chaos-ns"),
		Development:       r.Int()%2 == 0,
		StageDuration:     time.Duration(-10+r.Intn(200)) * time.Second,
		StageInterval:     time.Duration(-10+r.Intn(200)) * time.Second,
	})
}

// FromEnvironment loads Config from env vars.
func FromEnvironment() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, ErrParse
	}

	return cfg, nil
}
