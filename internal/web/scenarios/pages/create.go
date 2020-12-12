package pages

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"net/http"
)

func ScenarioCreationPage(rw http.ResponseWriter) {
	server.HTMLPage(rw, "web/html/scenarios/create.gohtml", nil)
}
