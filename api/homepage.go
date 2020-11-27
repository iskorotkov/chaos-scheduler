package api

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func Homepage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		server.ReturnHTMLPage(rw, "templates/html/home.html", nil)
	} else {
		server.MethodNotAvailable(rw, r)
	}
}
