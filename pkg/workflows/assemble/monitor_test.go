package assemble

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/generate"
	"math/rand"
	"testing"
	"testing/quick"
)

func Test_stageMonitor_Apply(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(0))
	f := func(monitor monitor, stage generate.Stage) bool {
		index := r.Intn(100)
		templates := monitor.Apply(stage, index)

		if len(templates) != 1 {
			t.Log("monitor extension must add exactly one template")
			return false
		}

		if templates[0].Name != fmt.Sprintf("monitor-%d", index+1) ||
			templates[0].Container.Name != "monitor" ||
			templates[0].Container.Image != monitor.image {
			t.Log("monitor template has incorrect values")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{Rand: r}); err != nil {
		t.Fatal(err)
	}
}
