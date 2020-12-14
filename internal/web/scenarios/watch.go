package scenarios

import (
	"github.com/go-chi/chi"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/watcher"
	"github.com/iskorotkov/chaos-scheduler/pkg/ws"
	"go.uber.org/zap"
	"net/http"
	"time"
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

	socket, err := ws.NewWebsocket(w, r, time.Hour*2, logger.Named("websocket"))
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "couldn't create websocket connection", http.StatusInternalServerError)
		return
	}

	m := watcher.NewMonitor(cfg.ServerURL, logger.Named("monitor"))

	events := make(chan *watcher.Event)
	clientLeft := make(chan struct{})

	go func() {
		if err := socket.WaitForClose(clientLeft); err != nil {
			logger.Error(err.Error())
		}
	}()

	go func() {
		// Read all remaining events
		defer func() {
			for {
				if event := <-events; event == nil {
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

	loop:
		for {
			select {
			case event := <-events:
				if event == nil {
					break
				}

				if err := socket.Write(event); err != nil {
					logger.Error(err.Error())
					break
				}
			case <-clientLeft:
				break loop
			}
		}

		logger.Info("all workflow events were processed")
	}()

	go func() {
		if err := m.Start(req.Name, req.Namespace, events); err != nil {
			logger.Error(err.Error())
		}

		logger.Info("all workflow events were sent")
	}()
}
