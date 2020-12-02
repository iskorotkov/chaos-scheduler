package importers

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
	"io/ioutil"
	"path"
	"strings"
)

type FolderImporter struct {
	Path string
}

func (f FolderImporter) Import() ([]scenarios.Template, error) {
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		logger.Error(err)
		return nil, FolderNotFoundError
	}

	actions := make([]scenarios.Template, 0)
	for _, file := range files {
		b, err := ioutil.ReadFile(path.Join(f.Path, file.Name()))
		if err != nil {
			logger.Error(err)
			return nil, FileError
		}

		filename := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))

		meta := Metadata{name: filename}
		actions = append(actions, item{meta, string(b)})
	}

	return actions, nil
}

func NewFolderImporter(path string) Importer {
	return FolderImporter{path}
}
