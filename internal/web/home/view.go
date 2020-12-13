package home

import (
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"go.uber.org/zap"
	"net/http"
)

func view(w http.ResponseWriter, r *http.Request, l *zap.SugaredLogger) {
	server.PageHandler("web/html/home/home.gohtml", nil, l)(w, r)
}
