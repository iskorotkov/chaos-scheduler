package execute

import (
	"context"
	"errors"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/templates"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

var (
	ConnectionError = errors.New("couldn't post scenario to executor server")
	ResponseError   = errors.New("executor server returned invalid status code")
)

func Execute(url string, wf templates.Workflow, logger *zap.SugaredLogger) (templates.Workflow, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Errorw(err.Error(),
			"url", url)
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
		logger.Errorw(err.Error(),
			"workflow", wf)
		return templates.Workflow{}, ResponseError
	}

	return templates.Workflow(*createdWf), nil
}