package config

import (
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"os"
	"time"
)

type Config struct {
	ServerURL         string
	StageMonitorImage string
	AppNS             string
	AppLabel          string
	ChaosNS           string
	IsInKubernetes    bool
	StageDuration     time.Duration
	StageInterval     time.Duration
}

func ParseConfigFromEnv() Config {
	logger.SetLevel(os.Getenv("LOGGING_LEVEL"))

	url := os.Getenv("SERVER_URL")
	if url == "" {
		logger.Critical(errors.New("executor host isn't set"))
	}

	stageMonitorImage := os.Getenv("STAGE_MONITOR_IMAGE")
	if stageMonitorImage == "" {
		logger.Warning("stage monitor image isn't specified; no stage monitor will be created")
	}

	appNS := os.Getenv("APP_NS")
	if appNS == "" {
		appNS = "chaos-app"
		logger.Warning(fmt.Sprintf("target namespace isn't set; using default value of %s", appNS))
	}

	appLabel := os.Getenv("APP_LABEL")
	if appLabel == "" {
		appLabel = "app"
		logger.Warning(fmt.Sprintf("app label isn't set; using default value of %s", appLabel))
	}

	chaosNS := os.Getenv("CHAOS_NS")
	if chaosNS == "" {
		chaosNS = "chaos-framework"
		logger.Warning(fmt.Sprintf("infrastructure namespace isn't set; using default value of %s", chaosNS))
	}

	isInKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	stageDuration, err := time.ParseDuration(os.Getenv("STAGE_DURATION"))
	if err != nil {
		stageDuration = time.Second * 30
		logger.Warning(fmt.Sprintf("stage duration isn't set; using default value of %s", stageDuration.String()))
	}

	stageInterval, err := time.ParseDuration(os.Getenv("STAGE_INTERVAL"))
	if err != nil {
		stageInterval = time.Second * 30
		logger.Warning(fmt.Sprintf("stage interval isn't set; using default value of %s", stageInterval.String()))
	}

	return Config{
		ServerURL:         url,
		StageMonitorImage: stageMonitorImage,
		AppNS:             appNS,
		AppLabel:          appLabel,
		ChaosNS:           chaosNS,
		IsInKubernetes:    isInKubernetes,
		StageDuration:     stageDuration,
		StageInterval:     stageInterval,
	}
}
