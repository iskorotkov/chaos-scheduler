package files

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

func ServeFolder(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		urlPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dir)))
		fs.ServeHTTP(w, r)
	}
}
