package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ConnectionError = errors.New("couldn't establish connection to gRPC server")
	RequestError    = errors.New("couldn't start watching updates")
	StreamError     = errors.New("couldn't read workflow update")
	SpecError       = errors.New("couldn't find step spec")
)

type Watcher struct {
	url    string
	logger *zap.SugaredLogger
}

func NewMonitor(url string, logger *zap.SugaredLogger) Watcher {
	return Watcher{url: url, logger: logger}
}

func (w Watcher) Start(ctx context.Context, name string, namespace string, output chan<- *Event) error {
	w.logger.Info("opening watcher gRPC connection")

	conn, err := grpc.Dial(w.url, grpc.WithInsecure())
	if err != nil {
		w.logger.Errorw(err.Error(),
			"url", w.url)
		return ConnectionError
	}

	defer func() {
		w.logger.Info("closing watcher gRPC connection")
		err := conn.Close()
		if err != nil {
			w.logger.Error(err.Error())
		}
	}()

	client := workflow.NewWorkflowServiceClient(conn)

	selector := fmt.Sprintf("metadata.name=%s", name)
	options := &v1.ListOptions{FieldSelector: selector}
	request := &workflow.WatchWorkflowsRequest{Namespace: namespace, ListOptions: options}
	service, err := client.WatchWorkflows(ctx, request)
	if err != nil {
		w.logger.Errorw(err.Error(),
			"selector", selector)
		return RequestError
	}

	defer close(output)

	for {
		msg, err := service.Recv()
		if err == io.EOF || ctx.Err() != nil {
			break
		}

		if err != nil {
			w.logger.Errorw(err.Error(),
				"selector", selector,
				"namespace", namespace)
			return StreamError
		}

		ev, err := newEvent(msg)
		if err != nil {
			w.logger.Error(err.Error())
			return err
		}

		output <- ev

		if ev.Phase != "Running" && ev.Phase != "Pending" {
			break
		}
	}

	return nil
}
