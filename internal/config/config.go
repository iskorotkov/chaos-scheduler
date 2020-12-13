package config

import (
	"errors"
	"github.com/caarlos0/env"
	"time"
)

var (
	ParseError = errors.New("couldn't parse config from env vars")
)

type Config struct {
	ServerURL         string        `env:"SERVER_URL"`
	StageMonitorImage string        `env:"STAGE_MONITOR_IMAGE"`
	AppNS             string        `env:"APP_NS"`
	AppLabel          string        `env:"APP_LABEL"`
	ChaosNS           string        `env:"CHAOS_NS"`
	Development       bool          `env:"DEVELOPMENT"`
	StageDuration     time.Duration `env:"STAGE_DURATION"`
	StageInterval     time.Duration `env:"STAGE_INTERVAL"`
}

func FromEnvironment() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, ParseError
	}

	return cfg, nil
}
