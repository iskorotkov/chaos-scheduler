package executors

import (
	"context"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type GRPCExecutor struct {
	url    string
	logger *zap.SugaredLogger
}

func NewGRPCExecutor(url string, logger *zap.SugaredLogger) GRPCExecutor {
	return GRPCExecutor{url: url, logger: logger}
}

func (g GRPCExecutor) Execute(wf templates.Workflow) (templates.Workflow, error) {
	conn, err := grpc.Dial(g.url, grpc.WithInsecure())
	if err != nil {
		g.logger.Errorw(err.Error(),
			"url", g.url)
		return templates.Workflow{}, ConnectionError
	}

	client := workflow.NewWorkflowServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	argoWf := v1alpha1.Workflow(wf)
	createdWf, err := client.CreateWorkflow(ctx, &workflow.WorkflowCreateRequest{
		Namespace: wf.Namespace,
		Workflow:  &argoWf,
	})
	if err != nil {
		g.logger.Errorw(err.Error(),
			"workflow", wf)
		return templates.Workflow{}, ResponseError
	}

	return templates.Workflow(*createdWf), nil
}
