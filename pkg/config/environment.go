package config

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"os"
)

type Config struct {
	ServerURL         string
	StageMonitorImage string
	AppNS             string
	ChaosNS           string
	IsInKubernetes    bool
}

func ParseConfigFromEnv() Config {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		logger.Critical(errors.New("executor host isn't set"))
	}

	stageMonitorImage := os.Getenv("STAGE_MONITOR_IMAGE")
	if stageMonitorImage == "" {
		logger.Warning("stage monitor image wasn't specified; no stage monitor will be created")
	}

	appNS := os.Getenv("APP_NS")
	if appNS == "" {
		appNS = "default"
		logger.Warning("target namespace wasn't set")
	}

	chaosNS := os.Getenv("CHAOS_NS")
	if chaosNS == "" {
		chaosNS = "default"
		logger.Warning("infrastructure namespace wasn't set")
	}

	isInKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	return Config{
		ServerURL:         url,
		StageMonitorImage: stageMonitorImage,
		AppNS:             appNS,
		ChaosNS:           chaosNS,
		IsInKubernetes:    isInKubernetes,
	}
}
