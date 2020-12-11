package config

import (
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"os"
)

type Config struct {
	ServerURL         string
	StageMonitorImage string
	AppNS             string
	AppLabel          string
	ChaosNS           string
	IsInKubernetes    bool
}

func ParseConfigFromEnv() Config {
	logger.SetLevel(os.Getenv("LOGGING_LEVEL"))

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

	appLabel := os.Getenv("APP_LABEL")
	if appNS == "" {
		appNS = "app"
		logger.Warning("target namespace wasn't set")
	}

	chaosNS := os.Getenv("CHAOS_NS")
	if chaosNS == "" {
		chaosNS = "default"
		logger.Warning("infrastructure namespace wasn't set")
	}

	isInKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	cfg := Config{
		ServerURL:         url,
		StageMonitorImage: stageMonitorImage,
		AppNS:             appNS,
		AppLabel:          appLabel,
		ChaosNS:           chaosNS,
		IsInKubernetes:    isInKubernetes,
	}

	logger.Info(fmt.Sprintf("Launching with config: %#v", cfg))

	return cfg
}
