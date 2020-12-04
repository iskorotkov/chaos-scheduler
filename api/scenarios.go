package api

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/exporters"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/scenarios"
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

func Scenarios(rw http.ResponseWriter, r *http.Request, cfg config.Config) {
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

func submissionStatusPage(rw http.ResponseWriter, r *http.Request, cfg config.Config) {
	form, err := parseForm(r)
	if err != nil {
		server.BadRequest(rw, err)
		return
	}

	workflow, err := generateWorkflow(rw, form, cfg)
	if err != nil {
		return
	}

	executor := executors.NewRestExecutor(cfg.ServerURL)
	err = executor.Execute(workflow)
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioExecutionError)
		return
	}

	server.HTMLPage(rw, "static/html/scenarios/submission-status.gohtml", nil)
}

func generateWorkflow(rw http.ResponseWriter, form Form, cfg config.Config) (string, error) {
	generator := scenarios.NewRoundRobinGenerator()
	exporter := exporters.NewJsonExporter()
	assembler := createAssembler(cfg)

	workflow, err := workflows.NewWorkflow(workflows.WorkflowParams{
		Generator: generator,
		Config: scenarios.ScenarioParams{
			Stages: form.Stages,
			Seed:   form.Seed,
		},
		Assembler: assembler,
		Exporter:  exporter,
	})
	if err != nil {
		if err == workflows.TemplatesImportError || err == workflows.WorkflowExportError {
			server.InternalError(rw, err)
		} else {
			server.BadRequest(rw, err)
		}

		return "", err
	}

	return workflow, nil
}

func createAssembler(cfg config.Config) assemblers.Assembler {
	return assemblers.NewModularAssembler(
		nil,
		[]extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.TargetNamespace), extensions.UseSuspend()},
		[]extensions.WorkflowExtension{extensions.UseSteps()},
	)
}

func scenarioCreationPage(rw http.ResponseWriter) {
	server.HTMLPage(rw, "static/html/scenarios/create.gohtml", nil)
}

func scenarioPreviewPage(rw http.ResponseWriter, cfg config.Config, form Form) {
	workflow, err := generateWorkflow(rw, form, cfg)
	if err != nil {
		return
	}

	server.HTMLPage(rw, "static/html/scenarios/preview.gohtml", struct {
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
