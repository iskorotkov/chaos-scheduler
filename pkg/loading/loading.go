package loading

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenario"
	"io/ioutil"
	"path"
)

var _ scenario.Failure = Failure{}

type Failure struct {
	name, yaml string
}

func (f Failure) Name() string {
	return f.name
}

func (f Failure) Yaml() string {
	return f.yaml
}

func Load(folder string) ([]Failure, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("folder doesn't exist: %v", err)
	}

	failures := make([]Failure, 0)
	for _, file := range files {
		b, err := ioutil.ReadFile(path.Join(folder, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("couldn't read file: %v", err)
		}

		content := string(b)
		failures = append(failures, Failure{file.Name(), content})
	}

	return failures, nil
}
