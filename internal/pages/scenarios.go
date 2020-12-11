package pages

import (
	"encoding/json"
	"errors"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/executors"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments/concrete"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generators"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"net/http"
	"strconv"
)

var (
	FormParseError         = errors.New("couldn't parse form data")
	ScenarioExecutionError = errors.New("couldn't execute scenario")
	MarshalError           = errors.New("couldn't marshall workflow to readable format")
	ScenarioParamsError    = errors.New("couldn't create scenario with given parameters")
	ScenarioGeneratorError = errors.New("couldn't generate scenario due to unknown reason")
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

	wf, err := generateWorkflow(form, cfg)
	if err != nil {
		if err == ScenarioParamsError {
			server.BadRequest(rw, err)
		} else {
			server.InternalError(rw, err)
		}

		return
	}

	executor := executors.NewGRPCExecutor(cfg.ServerURL)
	wf, err = executor.Execute(wf)
	if err != nil {
		logger.Error(err)
		server.InternalError(rw, ScenarioExecutionError)
		return
	}

	server.HTMLPage(rw, "web/html/scenarios/submission-status.gohtml", nil)
}

func generateWorkflow(form Form, cfg config.Config) (templates.Workflow, error) {
	presetList := experiments.List{
		PodPresets: []experiments.PodEnginePreset{
			concrete.PodDelete{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, Interval: 1, Force: false},
		},
		ContainerPresets: []experiments.ContainerEnginePreset{
			concrete.PodNetworkLatency{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, NetworkLatency: 300},
			concrete.PodNetworkLoss{Namespace: cfg.ChaosNS, AppNamespace: cfg.AppNS, LossPercentage: 100},
		},
	}

	extensionsList := extensions.List{
		ActionExtensions:   nil,
		StageExtensions:    []extensions.StageExtension{extensions.UseStageMonitor(cfg.StageMonitorImage, cfg.AppNS, cfg.StageInterval)},
		WorkflowExtensions: []extensions.WorkflowExtension{extensions.UseSteps()},
	}

	seeker, err := targets.NewSeeker(cfg.AppNS, cfg.AppLabel, cfg.IsInKubernetes)
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, ScenarioGeneratorError
	}

	g := generators.NewRoundRobinGenerator(presetList, seeker)
	s, err := g.Generate(generators.Params{Stages: form.Stages, Seed: form.Seed, StageDuration: cfg.StageDuration})
	if err != nil {
		logger.Error(err)

		if err == generators.NonPositiveStagesError || err == generators.TooManyStagesError {
			return templates.Workflow{}, ScenarioParamsError
		}

		return templates.Workflow{}, ScenarioGeneratorError
	}

	a := assemblers.NewModularAssembler(extensionsList)
	wf, err := a.Assemble(s)
	if err != nil {
		logger.Error(err)
		return templates.Workflow{}, ScenarioGeneratorError
	}

	return wf, nil
}

func scenarioCreationPage(rw http.ResponseWriter) {
	server.HTMLPage(rw, "web/html/scenarios/create.gohtml", nil)
}

func scenarioPreviewPage(rw http.ResponseWriter, cfg config.Config, form Form) {
	wf, err := generateWorkflow(form, cfg)
	if err != nil {
		if err == ScenarioParamsError {
			server.BadRequest(rw, err)
		} else {
			server.InternalError(rw, err)
		}

		return
	}

	marshaled, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		logger.Error(err)
		server.BadRequest(rw, MarshalError)
		return
	}

	server.HTMLPage(rw, "web/html/scenarios/preview.gohtml", struct {
		GeneratedWorkflow string
		Seed              int64
		Stages            int
	}{string(marshaled), form.Seed, form.Stages})
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
