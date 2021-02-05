package server

import (
	"go.uber.org/zap"
	"net/http"
)

func WithLogger(f func(w http.ResponseWriter, r *http.Request, l *zap.SugaredLogger), l *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, l)
	}
}
