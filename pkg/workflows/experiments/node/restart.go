package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type Restart struct {
	Namespace    string
	AppNamespace string
}

func (r Restart) Engine(target targets.Target, duration time.Duration) experiments.Engine {
	return r.Instantiate(target.Selector(), target.Node, duration)
}

func (r Restart) Info() experiments.Info {
	return experiments.Info{
		Name:   "node-restart",
		Lethal: true,
	}
}

func (r Restart) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        r.Info().Name,
		Namespace:   r.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: experiments.AppInfo{
			AppNS:    r.AppNamespace,
			AppLabel: label,
			AppKind:  "deployment",
		},
		Experiments: []experiments.Experiment{
			experiments.NewExperiment(experiments.ExperimentParams{
				Name: r.Info().Name,
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODE":          node,
					"SSH_USER":             "root",
					"TARGET_NODE_IP":       node,
					"REBOOT_COMMAND":       "sudo systemctl reboot",
				},
			}),
		},
	})
}
