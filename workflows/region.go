package workflows

import (
	"github.com/temporalio/cloud-samples-go/protogen/temporal/api/cloud/cloudservice/v1"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	// region workflow types
	GetRegionWorkflowType     = workflowPrefix + "get-region"
	GetAllRegionsWorkflowType = workflowPrefix + "get-regions"
)

type (
	RegionWorkflows interface {
		// Region Management Workflows
		GetRegion(ctx workflow.Context, in *cloudservice.GetRegionRequest) (*cloudservice.GetRegionResponse, error)
		GetAllRegions(ctx workflow.Context, in *cloudservice.GetRegionsRequest) (*cloudservice.GetRegionsResponse, error)
	}
)

func registerRegionWorkflows(w worker.Worker, wf RegionWorkflows) {
	for k, v := range map[string]any{
		GetRegionWorkflowType:     wf.GetRegion,
		GetAllRegionsWorkflowType: wf.GetAllRegions,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}
}

// Get a region
func (w *workflows) GetRegion(ctx workflow.Context, in *cloudservice.GetRegionRequest) (*cloudservice.GetRegionResponse, error) {
	return activities.GetRegion(withInfiniteRetryActivityOptions(ctx), in)
}

// Get multiple regions
func (w *workflows) GetAllRegions(ctx workflow.Context, in *cloudservice.GetRegionsRequest) (*cloudservice.GetRegionsResponse, error) {
	return activities.GetRegions(withInfiniteRetryActivityOptions(ctx), in)
}
