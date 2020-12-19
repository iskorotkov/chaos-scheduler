package advanced

type Scale int

type Severity int

type Cost float32

const (
	ScaleContainer Scale = iota
	ScalePod
	ScaleDeploymentPart
	ScaleDeployment
	ScaleNode
)

const (
	SeverityNonCritical Severity = iota
	SeverityCritical
	SeverityLethal
)

type Budget struct {
	MaxFailures int
	MaxPoints   Cost
}

type Modifiers struct {
	ByScale    map[Scale]Cost
	BySeverity map[Severity]Cost
}

func calculateCost(modifiers Modifiers, f Failure) Cost {
	severity, ok := modifiers.BySeverity[f.Severity]
	if !ok {
		severity = 1
	}

	scale, ok := modifiers.ByScale[f.Scale]
	if !ok {
		scale = 1
	}

	return severity * scale
}
