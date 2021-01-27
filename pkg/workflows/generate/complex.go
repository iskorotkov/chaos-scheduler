package generate

func addComplexFailures(params Params) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageFailures := params.Failures
		stageTargets := params.Targets

		actions := make([]Action, 0)

		points := params.Budget.MaxPoints
		retries := params.Retries
		for len(actions) < params.Budget.MaxFailures {
			failure := randomFailure(stageFailures, params.RNG)
			target := randomTarget(stageTargets, params.RNG)
			cost := calculateCost(params.Modifiers, failure)

			if cost <= points {
				points -= cost

				actions = append(actions, Action{
					Name:     failure.Name(),
					Severity: failure.Severity,
					Scale:    failure.Scale,
					Target:   target,
					Engine:   failure.Template.Instantiate(target, params.StageDuration),
				})
			} else {
				if retries <= 0 {
					break
				}

				retries--
			}
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
