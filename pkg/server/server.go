package server

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"html/template"
	"net/http"
)

func WithConfig(f func(http.ResponseWriter, *http.Request, Config), cfg Config) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, cfg)
	}
}

func HTMLPage(rw http.ResponseWriter, path string, data interface{}) {
	t, err := template.ParseFiles(path)
	if err != nil {
		logger.Error(err)
		return
	}

	err = t.Execute(rw, data)
	if err != nil {
		logger.Error(err)
		return
	}
}

func MethodNotAvailable(rw http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("method %s is not supported for path %s", r.Method, r.URL.Path)
	logger.Error(err)
	http.Error(rw, err.Error(), http.StatusBadRequest)
}

func BadRequest(rw http.ResponseWriter, err error) {
	logger.Error(err)
	http.Error(rw, err.Error(), http.StatusBadRequest)
}

func InternalError(rw http.ResponseWriter, err error) {
	logger.Error(err)
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
