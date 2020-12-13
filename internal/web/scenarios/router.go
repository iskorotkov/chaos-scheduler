package scenarios

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/preview", Preview)

	r.Get("/create", server.Page("web/html/scenarios/create.gohtml", nil))
	r.Post("/create", CreatePost)

	r.Get("/view/{namespace}/{name}", Status)

	return r
}
