package argo

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

func load(folder string) ([]Action, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("folder doesn't exist: %v", err)
	}

	actions := make([]Action, 0)
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

		filename := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		actions = append(actions, Action{filename, content})
	}

	return actions, nil
}
