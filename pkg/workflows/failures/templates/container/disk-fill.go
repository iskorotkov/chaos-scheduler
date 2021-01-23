package container

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures/templates"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type DiskFill struct {
	Namespace      string
	AppNamespace   string
	FillPercentage int
}

func (d DiskFill) Instantiate(target targets.Target, duration time.Duration) templates.Engine {
	return templates.NewEngine(templates.EngineParams{
		Name:        d.Name(),
		Namespace:   d.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: templates.AppInfo{
			AppNS:    d.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []templates.Experiment{
			templates.NewExperiment(templates.ExperimentParams{
				Name: d.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_CONTAINER":     target.MainContainer,
					"FILL_PERCENTAGE":      strconv.Itoa(d.FillPercentage),
				},
			}),
		},
	})
}

func (d DiskFill) Name() string {
	return "disk-fill"
}
