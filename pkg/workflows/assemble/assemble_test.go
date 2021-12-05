package assemble

import (
	"math/rand"
	"testing"
	"testing/quick"

	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"github.com/iskorotkov/metadata"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestModularAssembler_Assemble(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	hasStageWithZeroActions := func(stages []generate.Stage) bool {
		for _, s := range stages {
			if len(s.Steps) == 0 {
				return true
			}
		}
		return false
	}

	f := func(scenario generate.Scenario, ext ExtCollection) bool {
		wf, err := Assemble(scenario, ext)
		if err == ErrStages && len(scenario.Stages) == 0 {
			return true
		} else if err == ErrActions && hasStageWithZeroActions(scenario.Stages) {
			return true
		} else if err != nil {
			return false
		}

		if wf.Namespace == "" ||
			wf.GenerateName == "" ||
			wf.Spec.ServiceAccountName == "" ||
			wf.Spec.Entrypoint == "" {
			return false
		}

		for _, template := range wf.Spec.Templates {
			if template.Name == "" ||
				len(template.Metadata.Labels) != 0 ||
				len(template.Metadata.Annotations) != 4 {
				return false
			}

			// TODO: Do not use temporary ObjectMeta to unmarshal metadata
			var objectMeta = v1.ObjectMeta{
				Labels:      template.Metadata.Labels,
				Annotations: template.Metadata.Annotations,
			}

			var meta api.TemplateMetadata
			err := metadata.Unmarshal(objectMeta, &meta, api.Prefix)
			if err != nil {
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Error(err)
	}
}
