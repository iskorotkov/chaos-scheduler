package api

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/execution"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/output"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/scenario"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
	"strconv"
)

type Form struct {
	Seed   int64
	Stages int
}

func Scenarios(rw http.ResponseWriter, r *http.Request, cfg server.Config) {
	if r.Method == "GET" {
		form, err := parseForm(r)
		if err != nil {
			scenarioCreationPage(rw)
		} else {
			scenarioPreviewPage(rw, cfg, form)
		}
	} else if r.Method == "POST" {
		submissionStatusPage(rw, r, cfg)
	} else {
		server.MethodNotAvailable(rw, r)
	}
}

func submissionStatusPage(rw http.ResponseWriter, r *http.Request, cfg server.Config) {
	form, err := parseForm(r)
	if err != nil {
		server.InternalError(rw, fmt.Sprintf("couldn't parse form data: %v", err))
		return
	}

	s, err := scenario.NewScenario(scenario.Config{Path: cfg.TemplatesPath, Stages: form.Stages})
	if err != nil {
		server.InternalError(rw, fmt.Sprintf("couldn't create test scenario: %v", err))
		return
	}

	err = execution.ExecuteFromConfig(cfg.ServerURL, output.Config{TemplatePath: cfg.WorkflowTemplatePath, Scenario: s})
	if err != nil {
		server.InternalError(rw, fmt.Sprintf("couldn't execute generated scenario: %v", err))
		return
	}

	server.ReturnHTMLPage(rw, "templates/html/scenarios/submission-status.gohtml", nil)
}

func scenarioCreationPage(rw http.ResponseWriter) {
	server.ReturnHTMLPage(rw, "templates/html/scenarios/create.gohtml", nil)
}

func scenarioPreviewPage(rw http.ResponseWriter, cfg server.Config, form Form) {
	s, err := scenario.NewScenario(scenario.Config{Path: cfg.TemplatesPath, Stages: form.Stages, Seed: form.Seed})
	if err != nil {
		server.InternalError(rw, fmt.Sprintf("couldn't create test scenario: %v", err))
		return
	}

	out, err := output.GenerateFromConfig(output.Config{TemplatePath: cfg.WorkflowTemplatePath, Scenario: s})
	if err != nil {
		server.InternalError(rw, fmt.Sprintf("couldn't convert scenario to given format: %v", err))
		return
	}

	server.ReturnHTMLPage(rw, "templates/html/scenarios/preview.gohtml", struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{out, form.Seed, form.Stages})
}

func parseForm(r *http.Request) (Form, error) {
	err := r.ParseForm()
	if err != nil {
		return Form{}, fmt.Errorf("couldn't parse form data: %v", err)
	}

	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		return Form{}, fmt.Errorf("couldn't parse number of stages: %v", err)
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		return Form{}, fmt.Errorf("couldn't parse seed value: %v", err)
	}

	return Form{Seed: seed, Stages: int(stages)}, err
}
