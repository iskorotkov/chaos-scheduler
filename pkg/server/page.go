package server

import (
	"go.uber.org/zap"
	"html/template"
	"net/http"
)

func PageHandler(path string, data interface{}, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(path)
		if err != nil {
			logger.Errorw(err.Error(),
				"path", path)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			logger.Errorw(err.Error(),
				"path", path,
				"data", data)
			return
		}
	}
}
