package importers

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

var (
	FolderNotFoundError = errors.New("couldn't find specified folder")
	FileError           = errors.New("couldn't read template file")
)

var (
	_ scenarios.Template = item{}
)

type Metadata struct {
	name string
}

type item struct {
	metadata Metadata
	content  string
}

func (i item) Name() string {
	return i.metadata.name
}

func (i item) Template() string {
	return i.content
}

type Importer interface {
	Import() ([]scenarios.Template, error)
}
