package monitor

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
)

type WorkflowUpdate workflow.WorkflowWatchEvent

type Monitor struct {
	url    string
	logger *zap.SugaredLogger
}

func NewMonitor(url string, logger *zap.SugaredLogger) Monitor {
	return Monitor{url: url, logger: logger}
}

func (m Monitor) Start(name string, namespace string, output chan<- *WorkflowUpdate) error {
	m.logger.Info("opening monitor gRPC connection")

	conn, err := grpc.Dial(m.url, grpc.WithInsecure())
	if err != nil {
		m.logger.Errorw(err.Error(),
			"url", m.url)
		return ConnectionError
	}

	defer func() {
		m.logger.Info("closing monitor gRPC connection")
		err := conn.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}()

	client := workflow.NewWorkflowServiceClient(conn)

	selector := fmt.Sprintf("metadata.name=%s", name)
	options := &v1.ListOptions{FieldSelector: selector}
	request := &workflow.WatchWorkflowsRequest{Namespace: namespace, ListOptions: options}
	service, err := client.WatchWorkflows(context.Background(), request)
	if err != nil {
		m.logger.Errorw(err.Error(),
			"selector", selector)
		return RequestError
	}

	defer close(output)

	for {
		event, err := service.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			m.logger.Errorw(err.Error(),
				"event", event,
				"selector", selector,
				"namespace", namespace)
			return StreamError
		}

		output <- (*WorkflowUpdate)(event)
	}

	return nil
}
