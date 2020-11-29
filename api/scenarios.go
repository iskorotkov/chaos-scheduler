package api

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/executors"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/exporters"
	"github.com/iskorotkov/chaos-scheduler/pkg/argov2/importers"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
	"strconv"
)

var (
	FormParseError         = errors.New("couldn't parse form data")
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

	importer := importers.NewFolderImporter(cfg.TemplatesPath)
	generator := scenarios.NewRoundRobinGenerator()
	assembler := assemblers.NewSimpleAssembler(cfg.WorkflowTemplatePath)
	exporter := exporters.NewJsonExporter()

	workflow, err := argov2.NewWorkflow(argov2.Config{
		Importer:  importer,
		Generator: generator,
		Config: scenarios.Config{
			Stages: form.Stages,
			Seed:   form.Seed,
		},
		Assembler: assembler,
		Exporter:  exporter,
	})
	if err != nil {
		logger.Error(err)
		if err == argov2.TemplatesImportError || err == argov2.WorkflowExportError {
			server.InternalError(rw, err)
		} else {
			server.BadRequest(rw, err)
		}

		return
	}

	executor := executors.NewRestExecutor(cfg.ServerURL)
	err = executor.Execute(workflow)
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
	importer := importers.NewFolderImporter(cfg.TemplatesPath)
	generator := scenarios.NewRoundRobinGenerator()
	assembler := assemblers.NewSimpleAssembler(cfg.WorkflowTemplatePath)
	exporter := exporters.NewJsonExporter()

	workflow, err := argov2.NewWorkflow(argov2.Config{
		Importer:  importer,
		Generator: generator,
		Config: scenarios.Config{
			Stages: form.Stages,
			Seed:   form.Seed,
		},
		Assembler: assembler,
		Exporter:  exporter,
	})
	if err != nil {
		logger.Error(err)
		if err == argov2.TemplatesImportError || err == argov2.WorkflowExportError {
			server.InternalError(rw, err)
		} else {
			server.BadRequest(rw, err)
		}

		return
	}

	server.HTMLPage(rw, "templates/html/scenarios/preview.gohtml", struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{workflow, form.Seed, form.Stages})
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
