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
	ParseError = errors.New("couldn't parse config from env vars")
)

type Config struct {
	ArgoServer        string        `env:"ARGO_SERVER"`
	StageMonitorImage string        `env:"STAGE_MONITOR_IMAGE"`
	AppNS             string        `env:"APP_NS"`
	AppLabel          string        `env:"APP_LABEL"`
	ChaosNS           string        `env:"CHAOS_NS"`
	Development       bool          `env:"DEVELOPMENT"`
	StageDuration     time.Duration `env:"STAGE_DURATION"`
	StageInterval     time.Duration `env:"STAGE_INTERVAL"`
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

func FromEnvironment() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, ParseError
	}

	return cfg, nil
}
