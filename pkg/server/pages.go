package server

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"html/template"
	"net/http"
)

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
