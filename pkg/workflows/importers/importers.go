package importers

import (
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenarios"
)

var (
	FolderNotFoundError = errors.New("couldn't find specified folder")
	FileError           = errors.New("couldn't read template file")
)

type Importer interface {
	Import() ([]scenarios.TemplatedAction, error)
}
