package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/experiments"
	"strconv"
	"time"
)

type Restart struct {
	Namespace    string
	AppNamespace string
}

func (r Restart) Type() experiments.ExperimentType {
	return "node-restart"
}

func (r Restart) Info() experiments.Info {
	return experiments.Info{Lethal: true}
}

func (r Restart) Instantiate(label string, node string, duration time.Duration) experiments.Engine {
	return experiments.NewEngine(experiments.EngineParams{
		Name:        string(r.Type()),
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
				Type: r.Type(),
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
