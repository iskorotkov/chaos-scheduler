// Package argo allows to launch generated workflows using Argo server.
package argo

import (
	"context"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/assemble"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/execute"
	"go.uber.org/zap"
)

type argo struct {
	client apiclient.Client
	logger *zap.SugaredLogger
}

func NewExecutor(url string, logger *zap.SugaredLogger) (execute.Executor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, apiClient, err := apiclient.NewClientFromOpts(apiclient.Opts{
		ArgoServerOpts: apiclient.ArgoServerOpts{
			URL:                url,
			Path:               "",
			Secure:             true,
			InsecureSkipVerify: true,
			HTTP1:              false,
		},
		InstanceID: "",
		AuthSupplier: func() string {
			return ""
		},
		ClientConfigSupplier: nil,
		Offline:              false,
		Context:              ctx,
	})
	if err != nil {
		logger.Errorw(err.Error(),
			"url", url)
		return nil, execute.ErrConnection
	}

	return &argo{
		client: apiClient,
		logger: logger,
	}, nil
}

// Execute workflow using Argo server.
func (a argo) Execute(wf assemble.Workflow) (assemble.Workflow, error) {
	serviceClient := a.client.NewWorkflowServiceClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	argoWf := v1alpha1.Workflow(wf)
	createdWf, err := serviceClient.CreateWorkflow(ctx, &workflow.WorkflowCreateRequest{
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
