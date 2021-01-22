package assemblers

import (
	api "github.com/iskorotkov/chaos-scheduler/api/metadata"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemblers/extensions"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generator"
	"github.com/iskorotkov/metadata"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func TestModularAssembler_Assemble(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	a := ModularAssembler{
		Extensions: extensions.List{
			ActionExtensions: []extensions.ActionExtension{
				// No action extensions implemented
			},
			StageExtensions: []extensions.StageExtension{
				extensions.UseSuspend(),
				extensions.UseStageMonitor("stage-monitor", "target-ns", time.Duration(r.Intn(60)), &zap.SugaredLogger{}),
			},
			WorkflowExtensions: []extensions.WorkflowExtension{
				extensions.UseSteps(),
			},
		},
		logger: &zap.SugaredLogger{},
	}

	hasStageWithZeroActions := func(stages []generator.Stage) bool {
		for _, s := range stages {
			if len(s.Actions) == 0 {
				return true
			}
		}
		return false
	}

	f := func(scenario generator.Scenario) bool {
		wf, err := a.Assemble(scenario)
		if err == StagesError && len(scenario.Stages) == 0 {
			return true
		} else if err == ActionsError && hasStageWithZeroActions(scenario.Stages) {
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
