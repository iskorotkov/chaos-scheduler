// Package server helps to configure requests handling.
package server

import (
	"go.uber.org/zap"
	"net/http"
)

// WithLogger allows to pass logger into a request handler.
func WithLogger(f func(w http.ResponseWriter, r *http.Request, l *zap.SugaredLogger), l *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, l)
	}
}
