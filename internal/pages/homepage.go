package pages

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Homepage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		server.HTMLPage(rw, "web/html/home.gohtml", nil)
	} else {
		server.MethodNotAvailable(rw, r)
	}
}
