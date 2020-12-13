package server

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"html/template"
	"net/http"
)

func Page(path string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(path)
		if err != nil {
			logger.Error(err)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			logger.Error(err)
			return
		}
	}
}
