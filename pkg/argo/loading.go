package argo

import (
	"fmt"
	"io/ioutil"
	"path"
)

type failure struct {
	name string
	yaml string
}

func Load(folder string) ([]failure, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("folder doesn't exist: %v", err)
	}

	failures := make([]failure, 0)
	for _, file := range files {
		b, err := ioutil.ReadFile(path.Join(folder, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("couldn't read file: %v", err)
		}

		content := string(b)
		failures = append(failures, failure{file.Name(), content})
	}

	return failures, nil
}
