package server

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"net/http"
)

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
