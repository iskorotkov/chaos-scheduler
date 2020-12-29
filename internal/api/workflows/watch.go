package workflows

import (
	"context"
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

	socket, err := ws.NewWebsocket(w, r, logger.Named("websocket"))
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "couldn't create websocket connection", http.StatusInternalServerError)
		return
	}

	m := watcher.NewMonitor(cfg.ArgoServer, logger.Named("monitor"))

	events := make(chan *watcher.Event)

	monitorCtx, monitorCancel := context.WithTimeout(context.Background(), time.Hour)
	wsCtx, wsCancel := context.WithTimeout(context.Background(), time.Hour)

	go func() {
		defer readRemainingEvents(events)
		defer closeWebsocket(socket, logger)

		defer monitorCancel()
		defer wsCancel()

		sendEvents(wsCtx, socket, events, logger)
	}()

	go func() {
		defer monitorCancel()
		readEvents(monitorCtx, m, req, events, logger)
	}()
}

func sendEvents(ctx context.Context, socket ws.Websocket, events chan *watcher.Event, logger *zap.SugaredLogger) {
	for {
		select {
		case event := <-events:
			if event == nil {
				return
			}

			if err := socket.Write(ctx, event); err != nil && err != ws.DeadlineExceededError && err != ws.ContextCancelledError {
				logger.Error(err.Error())
				return
			}
		case <-socket.Closed:
			return
		}
	}
}

func closeWebsocket(socket ws.Websocket, logger *zap.SugaredLogger) {
	if err := socket.Close(); err != nil && err != ws.DeadlineExceededError {
		logger.Error(err.Error())
	}

	logger.Info("all workflow events were processed")
}

func readRemainingEvents(events <-chan *watcher.Event) {
	for {
		if event := <-events; event == nil {
			break
		}
	}
}

func readEvents(ctx context.Context, m watcher.Watcher, req request, events chan<- *watcher.Event, logger *zap.SugaredLogger) {
	if err := m.Start(ctx, req.Name, req.Namespace, events); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("all workflow events were sent")
}
