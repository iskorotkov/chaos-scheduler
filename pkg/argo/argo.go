// Package argo allows to launch generated workflows using Argo server.
package argo

import (
	"context"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type argo struct {
	conn   *grpc.ClientConn
	logger *zap.SugaredLogger
}

func NewExecutor(url string, logger *zap.SugaredLogger) (execute.Executor, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Errorw(err.Error(),
			"url", url)
		return nil, execute.ErrConnection
	}

	return &argo{
		conn:   conn,
		logger: logger,
	}, nil
}

// Execute workflow using Argo server.
func (a argo) Execute(wf assemble.Workflow) (assemble.Workflow, error) {
	client := workflow.NewWorkflowServiceClient(a.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	argoWf := v1alpha1.Workflow(wf)
	createdWf, err := client.CreateWorkflow(ctx, &workflow.WorkflowCreateRequest{
		Namespace: wf.Namespace,
		Workflow:  &argoWf,
	})
	if err != nil {
		a.logger.Errorw(err.Error(),
			"workflow", wf)
		return assemble.Workflow{}, execute.ErrResponse
	}

	return assemble.Workflow(*createdWf), nil
}
