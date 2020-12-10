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
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/presets/concrete"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"net/http"
	"strconv"
)

var (
	FormParseError         = errors.New("couldn't parse form data")
	ScenarioExecutionError = errors.New("couldn't execute scenario")
	SeekerCreationFailed   = errors.New("couldn't create seeker instance")
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
	appNS := cfg.AppNS
	chaosNS := cfg.ChaosNS

	presetList := presets.List{
		PodPresets: []presets.PodEnginePreset{
			concrete.PodDelete{Namespace: chaosNS, AppNamespace: appNS, Duration: 60, Interval: 5, Force: false},
		},
		ContainerPresets: []presets.ContainerEnginePreset{
			concrete.PodNetworkLatency{Namespace: chaosNS, AppNamespace: appNS, NetworkLatency: 300},
			concrete.PodNetworkLoss{Namespace: chaosNS, AppNamespace: appNS, LossPercentage: 100},
		},
	}

	extensionsList := extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS)},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}

	seeker, err := targets.NewSeeker(appNS, cfg.AppLabel, cfg.IsInKubernetes)
	if err != nil {
		logger.Error(err)
		return "", SeekerCreationFailed
	}

	workflow, err := workflows.NewWorkflow(
		generators.NewRoundRobinGenerator(presetList, seeker),
		assemblers.NewModularAssembler(extensionsList),
		exporters.NewJsonExporter(),
		generators.Params{Stages: form.Stages, Seed: form.Seed},
	)
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
