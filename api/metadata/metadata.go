// Package metadata describes Kubernetes metadata used on generated templates.
package metadata

type (
	// Version describes what fields are available.
	Version string
	// Type describes template purpose.
	Type string
	// Severity describes amount of damage the template will deal to single target.
	Severity string
	// Scale describes area of damage of the template.
	Scale string
)

const (
	// Prefix used to avoid collisions of Kubernetes labels and annotations.
	Prefix = "chaosframework.com"

	// VersionV1 is the first version of the metadata.
	VersionV1 = Version("v1")

	// TypeFailure describes failures.
	TypeFailure = Type("failure")
	// TypeUtility describes utility templates not causing any damage.
	TypeUtility = Type("utility")

	// SeverityHarmless describes templates that do no harm to the targets (or don't have targets at all).
	SeverityHarmless = Severity("harmless")
	// SeverityLight describes templates that do light damage to targets.
	SeverityLight = Severity("light")
	// SeveritySevere describes templates that do severe damage to targets.
	SeveritySevere = Severity("severe")
	// SeverityCritical describes templates that destroy/restart targets.
	SeverityCritical = Severity("critical")

	// ScaleContainer describes templates that target single container.
	ScaleContainer = Scale("container")
	// ScalePod describes templates that target single pod.
	ScalePod = Scale("pod")
	// ScaleDeploymentPart describes templates that target part of the deployment.
	ScaleDeploymentPart = Scale("deployment part")
	// ScaleDeployment describes templates that target entire deployment.
	ScaleDeployment = Scale("deployment")
	// ScaleNode describes templates that target entire node.
	ScaleNode = Scale("node")
	// ScaleCluster describes templates that target entire cluster.
	ScaleCluster = Scale("cluster")
)

// TemplateMetadata describes additional info about template.
type TemplateMetadata struct {
	// Version describes a version of template metadata.
	Version Version `annotation:"version"`
	// Type describes a type of template.
	Type Type `annotation:"type"`
	// Severity describes a severity of the template.
	Severity Severity `annotation:"severity"`
	// Scale describes a scale of the template.
	Scale Scale `annotation:"scale"`
}
