package importers

import "errors"

var (
	FolderNotFoundError = errors.New("couldn't find specified folder")
	FileError           = errors.New("couldn't read template file")
)

type Metadata struct {
	Name        string
	Labels      []string
	Annotations []string
}

type Item struct {
	Metadata Metadata
	Content  string
}

type Importer interface {
	Import() ([]Item, error)
}
