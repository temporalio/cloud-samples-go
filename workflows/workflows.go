package workflows

import (
	"github.com/temporalio/cloud-samples-go/client/api"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
	"go.temporal.io/sdk/worker"
)

//go:generate mockgen -source workflows.go -destination workflows_mock.go -package workflow

const (
	workflowPrefix = "tmprlcloud-wf."
)

type (
	Workflows interface {
		UserWorkflows
		NamespaceWorkflows
		RegionWorkflows
		AsyncOperationWorkflows
	}
	workflows struct{}
)

func NewWorkflows() Workflows {
	return &workflows{}
}

func NewActivities(client *api.Client) *activities.Activities {
	return activities.NewActivities(client)
}

func Register(w worker.Worker, wf Workflows, a *activities.Activities) {
	// Register the workflows that we want to be able to use.
	registerUserWorkflows(w, wf)
	registerNamespaceWorkflows(w, wf)
	registerRegionWorkflows(w, wf)
	registerAsyncOperationWorkflows(w, wf)

	// Register the activities that the workflows will use.
	activities.Register(w, a)
}
