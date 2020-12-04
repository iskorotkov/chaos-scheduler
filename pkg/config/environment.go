package config

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"os"
)

type Config struct {
	ServerURL            string
	WorkflowTemplatePath string
	StageMonitorImage    string
	TargetNamespace      string
	IsInKubernetes       bool
}

func ParseConfigFromEnv() Config {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		logger.Critical(errors.New("executor host isn't set"))
	}

	workflowTemplate := os.Getenv("WORKFLOW_TEMPLATE_PATH")
	if workflowTemplate == "" {
		logger.Critical(errors.New("path to workflow template isn't set"))
	}

	stageMonitorImage := os.Getenv("STAGE_MONITOR_IMAGE")
	if stageMonitorImage == "" {
		logger.Warning("stage monitor image wasn't specified; no stage monitor will be created")
	}

	targetNs := os.Getenv("TARGET_NAMESPACE")
	if targetNs == "" {
		targetNs = "default"
		logger.Warning("target namespace wasn't set")
	}

	isInKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	return Config{
		ServerURL:            url,
		WorkflowTemplatePath: workflowTemplate,
		StageMonitorImage:    stageMonitorImage,
		TargetNamespace:      targetNs,
		IsInKubernetes:       isInKubernetes,
	}
}
