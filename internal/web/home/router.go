package home

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", server.Page("web/html/home/home.gohtml", nil))
	return r
}
