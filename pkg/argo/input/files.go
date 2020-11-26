package input

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

type Template struct {
	Filename string
	Yaml     string
}

func Load(folder string) ([]Template, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("folder doesn't exist: %v", err)
	}

	actions := make([]Template, 0)
	for _, file := range files {
		b, err := ioutil.ReadFile(path.Join(folder, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("couldn't read file: %v", err)
		}

		content := string(b)
		filename := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))

		actions = append(actions, Template{filename, content})
	}

	return actions, nil
}
