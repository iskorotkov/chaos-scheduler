package server

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/config"
	"net/http"
)

func WithConfig(f func(http.ResponseWriter, *http.Request, config.Config), cfg config.Config) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, cfg)
	}
}
