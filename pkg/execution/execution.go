package execution

import "github.com/iskorotkov/chaos-scheduler/pkg/scenario"

type Executor interface {
	Execute(scenario scenario.Scenario) error
}
