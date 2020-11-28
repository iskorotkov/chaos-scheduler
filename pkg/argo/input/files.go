package input

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"io/ioutil"
	"path"
	"strings"
)

type Template struct {
	Filename string
	Yaml     string
}

var (
	FolderNotFoundError = errors.New("couldn't find specified folder")
	FileError           = errors.New("couldn't read template file")
)

func Load(folder string) ([]Template, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		logger.Error(err)
		return nil, FolderNotFoundError
	}

	actions := make([]Template, 0)
	for _, file := range files {
		b, err := ioutil.ReadFile(path.Join(folder, file.Name()))
		if err != nil {
			logger.Error(err)
			return nil, FileError
		}

		content := string(b)
		filename := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))

		actions = append(actions, Template{filename, content})
	}

	return actions, nil
}
