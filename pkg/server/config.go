package server

import (
	"log"
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
		log.Fatalf("executor host isn't set")
	}

	templates := os.Getenv("TEMPLATES_PATH")
	if templates == "" {
		log.Fatalf("path to tempaltes isn't set")
	}

	workflowTemplate := os.Getenv("WORKFLOW_TEMPLATE_PATH")
	if workflowTemplate == "" {
		log.Fatalf("path to workflow template isn't set")
	}

	return Config{ServerURL: url, TemplatesPath: templates, WorkflowTemplatePath: workflowTemplate}
}
