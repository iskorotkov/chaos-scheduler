package server

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"os"
)

type Config struct {
	ServerURL            string
	TemplatesPath        string
	WorkflowTemplatePath string
}

func ParseConfig() Config {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		logger.Critical(errors.New("executor host isn't set"))
	}

	templates := os.Getenv("TEMPLATES_PATH")
	if templates == "" {
		logger.Critical(errors.New("path to templates isn't set"))
	}

	workflowTemplate := os.Getenv("WORKFLOW_TEMPLATE_PATH")
	if workflowTemplate == "" {
		logger.Critical(errors.New("path to workflow template isn't set"))
	}

	return Config{ServerURL: url, TemplatesPath: templates, WorkflowTemplatePath: workflowTemplate}
}
