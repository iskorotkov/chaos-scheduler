package api

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/execution"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/output"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo/scenario"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
	"strconv"
)

var (
	FormParseError         = errors.New("couldn't parse form data")
	ScenarioCreationError  = errors.New("couldn't create test scenario")
	ScenarioFormatError    = errors.New("couldn't format test scenario")
	ScenarioExecutionError = errors.New("couldn't execute scenario")
)

type Form struct {
	Seed   int64
	Stages int
}

func Scenarios(rw http.ResponseWriter, r *http.Request, cfg server.Config) {
	if r.Method == "GET" {
		err := r.ParseForm()
		if err != nil {
			logger.Error(err)
			server.BadRequest(rw, FormParseError)
		}

		if len(r.Form) == 0 {
			scenarioCreationPage(rw)
		} else {
			form, err := parseForm(r)
			if err != nil {
				server.BadRequest(rw, err)
			}

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
		server.BadRequest(rw, err)
		return
	}

	s, err := scenario.NewScenario(scenario.Config{Path: cfg.TemplatesPath, Stages: form.Stages})
	if err != nil {
		logger.Error(err)
		server.BadRequest(rw, ScenarioCreationError)
		return
	}

	err = execution.ExecuteFromConfig(cfg.ServerURL, output.Config{TemplatePath: cfg.WorkflowTemplatePath, Scenario: s})
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioExecutionError)
		return
	}

	server.HTMLPage(rw, "templates/html/scenarios/submission-status.gohtml", nil)
}

func scenarioCreationPage(rw http.ResponseWriter) {
	server.HTMLPage(rw, "templates/html/scenarios/create.gohtml", nil)
}

func scenarioPreviewPage(rw http.ResponseWriter, cfg server.Config, form Form) {
	s, err := scenario.NewScenario(scenario.Config{Path: cfg.TemplatesPath, Stages: form.Stages, Seed: form.Seed})
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioCreationError)
		return
	}

	out, err := output.GenerateFromConfig(output.Config{TemplatePath: cfg.WorkflowTemplatePath, Scenario: s})
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioFormatError)
		return
	}

	server.HTMLPage(rw, "templates/html/scenarios/preview.gohtml", struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{out, form.Seed, form.Stages})
}

func parseForm(r *http.Request) (Form, error) {
	err := r.ParseForm()
	if err != nil {
		logger.Error(err)
		return Form{}, FormParseError
	}

	stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
	if err != nil {
		logger.Error(err)
		return Form{}, FormParseError
	}

	seed, err := strconv.ParseInt(r.FormValue("seed"), 10, 64)
	if err != nil {
		logger.Error(err)
		return Form{}, FormParseError
	}

	return Form{Seed: seed, Stages: int(stages)}, err
}
