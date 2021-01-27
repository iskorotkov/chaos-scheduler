package generate

func addCascadeFailures(params Params) []Stage {
	phaseFailures := params.Failures

	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		stageTargets := params.Targets

		actions := make([]Action, 0)
		points := params.Budget.MaxPoints

		failure := randomFailure(phaseFailures, params.RNG)
		cost := calculateCost(params.Modifiers, failure)

		for i := 0; i < params.Retries; i++ {
			if cost*2 <= points {
				break
			}

			failure := randomFailure(phaseFailures, params.RNG)
			cost = calculateCost(params.Modifiers, failure)
		}

		for len(actions) < params.Budget.MaxFailures {
			if len(stageTargets) == 0 {
				break
			}

			target := randomTarget(stageTargets, params.RNG)

			actions = append(actions, Action{
				Name:     failure.Name(),
				Severity: failure.Severity,
				Scale:    failure.Scale,
				Target:   target,
				Engine:   failure.Template.Instantiate(target, params.StageDuration),
			})

			points -= cost
			if cost > points {
				break
			}
		}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
