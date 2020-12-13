package scenarios

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/monitor"
	"github.com/iskorotkov/chaos-scheduler/pkg/ws"
	"go.uber.org/zap"
	"net/http"
)

type request struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func watchWS(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	req := request{
		Name:      chi.URLParam(r, "name"),
		Namespace: chi.URLParam(r, "namespace"),
	}

	logger.Infow("get request params from url",
		"request", req)

	entry := r.Context().Value("config")
	cfg, ok := entry.(*config.Config)
	if !ok {
		msg := "couldn't get config for request context"
		logger.Errorw(msg,
			"config", "entry")
		http.Error(w, msg, http.StatusInternalServerError)
	}

	socket, err := ws.NewWebsocket(w, r, logger.Named("websocket"))
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "couldn't create websocket connection", http.StatusInternalServerError)
		return
	}

	m := monitor.NewMonitor(cfg.ServerURL, logger.Named("monitor"))

	updates := make(chan *monitor.WorkflowUpdate)

	go func() {
		// Read all remaining events
		defer func() {
			for {
				if update := <-updates; update == nil {
					break
				}
			}
		}()

		// Close socket
		defer func() {
			if err := socket.Close(); err != nil {
				logger.Error(err.Error())
			}
		}()

		for {
			update := <-updates
			if update == nil {
				break
			}

			if err := socket.Write(update); err != nil {
				logger.Error(err.Error())
				break
			}
		}

		logger.Info("all workflow updates were processed")
	}()

	go func() {
		if err := m.Start(req.Name, req.Namespace, updates); err != nil {
			logger.Error(err.Error())
		}

		logger.Info("all workflow updates were sent")
	}()
}
