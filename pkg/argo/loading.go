package argo

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type YamlContent map[interface{}]interface{}

type failure struct {
	name string
	yaml YamlContent
}

func load(folder string) ([]failure, error) {
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

		content := make(YamlContent)
		err = yaml.Unmarshal(b, content)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshall failure definition from yaml: %v", err)
		}

		failures = append(failures, failure{file.Name(), content})
	}

	return failures, nil
}
