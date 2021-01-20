package node

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"strconv"
	"time"
)

type Restart struct {
	Namespace    string
	AppNamespace string
}

func (r Restart) Instantiate(target targets.Target, duration time.Duration) failures.Engine {
	return failures.NewEngine(failures.EngineParams{
		Name:        r.Name(),
		Namespace:   r.Namespace,
		Labels:      nil,
		Annotations: nil,
		AppInfo: failures.AppInfo{
			AppNS:    r.AppNamespace,
			AppLabel: target.AppLabel,
			AppKind:  "deployment",
		},
		Experiments: []failures.Experiment{
			failures.NewExperiment(failures.ExperimentParams{
				Name: r.Name(),
				Env: map[string]string{
					"TOTAL_CHAOS_DURATION": strconv.Itoa(int(duration.Seconds())),
					"TARGET_NODE":          target.Node,
					"SSH_USER":             "root",
					"TARGET_NODE_IP":       target.Node,
					"REBOOT_COMMAND":       "sudo systemctl reboot",
				},
			}),
		},
	})
}

func (r Restart) Name() string {
	return "node-restart"
}
