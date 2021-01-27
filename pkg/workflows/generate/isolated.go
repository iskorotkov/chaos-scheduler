package generate

func addIsolatedFailures(params Params) []Stage {
	stages := make([]Stage, 0)

	for i := 0; i < params.Stages; i++ {
		failure := randomFailure(params.Failures, params.RNG)
		target := randomTarget(params.Targets, params.RNG)

		actions := []Action{{
			Name:     failure.Name(),
			Severity: failure.Severity,
			Scale:    failure.Scale,
			Target:   target,
			Engine:   failure.Template.Instantiate(target, params.StageDuration),
		}}

		stages = append(stages, Stage{
			Actions:  actions,
			Duration: params.StageDuration,
		})
	}

	return stages
}
