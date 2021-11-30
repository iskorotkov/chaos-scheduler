package handlers

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/failures"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
)

func mergeTargets(ts []targets.Target, ids []string) []targets.Target {
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	var res []targets.Target
	for _, t := range ts {
		if idMap[t.ID()] {
			res = append(res, t)
		}
	}

	return res
}

func mergeFailures(fs []failures.Failure, ids []string) []failures.Failure {
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	var res []failures.Failure
	for _, f := range fs {
		if idMap[f.ID()] {
			res = append(res, f)
		}
	}

	return res
}
