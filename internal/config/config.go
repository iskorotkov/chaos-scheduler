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
	// ArgoServer is an address of Argo used for executing generated workflows.
	ArgoServer string `env:"ARGO_SERVER"`
	// StageMonitorImage is a Docker image used for monitoring the state of the system under test.
	StageMonitorImage string `env:"STAGE_MONITOR_IMAGE"`
	// AppNS is a namespace of the system under test.
	AppNS string `env:"APP_NS"`
	// AppLabel is a metadata label used for target selection.
	AppLabel string `env:"APP_LABEL"`
	// ChaosNS is a namespace where workflows will be created.
	ChaosNS string `env:"CHAOS_NS"`
	// Development describes whether is in development or not.
	Development bool `env:"DEVELOPMENT"`
	// StageDuration describes a duration of each test stage.
	StageDuration time.Duration `env:"STAGE_DURATION"`
	// StageInterval describes an interval between test stages.
	StageInterval time.Duration `env:"STAGE_INTERVAL"`
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
