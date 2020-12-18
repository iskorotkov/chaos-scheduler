package generator

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
	MaxExperiments int
	MaxPoints      Cost
}

type Modifiers struct {
	ByScale    map[Scale]Cost
	BySeverity map[Severity]Cost
}
