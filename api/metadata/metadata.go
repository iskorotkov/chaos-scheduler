package metadata

type Version string
type Type string
type Utility string
type Severity string
type Scale string

const (
	Prefix = "chaosframework.com"

	VersionV1 = Version("v1")

	TypeFailure = Type("failure")
	TypeUtility = Type("utility")

	SeverityNonCritical = Severity("non critical")
	SeverityCritical    = Severity("critical")
	SeverityLethal      = Severity("lethal")

	ScaleContainer      = Scale("container")
	ScalePod            = Scale("pod")
	ScaleDeploymentPart = Scale("deployment part")
	ScaleDeployment     = Scale("deployment")
	ScaleNode           = Scale("node")
)

type TemplateMetadata struct {
	Version  Version  `annotation:"version"`
	Type     Type     `annotation:"type"`
	Severity Severity `annotation:"severity"`
	Scale    Scale    `annotation:"scale"`
}
