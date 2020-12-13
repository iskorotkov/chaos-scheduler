package scenarios

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func createPage(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	server.PageHandler("web/html/scenarios/create.gohtml", nil, logger.Named("page"))(w, r)
}
