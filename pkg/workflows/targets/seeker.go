package targets

import "errors"

var (
	ClientsetError = errors.New("couldn't create clientset")
	FetchError     = errors.New("couldn't fetch info from Kubernetes")
)

type Seeker interface {
	Targets() ([]Target, error)
}
