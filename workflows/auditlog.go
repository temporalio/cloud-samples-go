package workflows

import (
	"github.com/temporalio/cloud-samples-go/workflows/activities"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	// audit workflow types
	GetAuditLogsWorkflowType = workflowPrefix + "get-audit-logs"
)

type (
	AuditWorkflows interface {
		//Audit Workflows
		GetAuditLogs(ctx workflow.Context, in *cloudservice.GetAuditLogsRequest) (*cloudservice.GetAuditLogsResponse, error)
	}
)

func registerAuditWorkflows(w worker.Worker, wf AuditWorkflows) {
	for k, v := range map[string]any{
		GetAuditLogsWorkflowType: wf.GetAuditLogs,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}
}

// GetAuditLogs is a workflow that retrieves audit logs from Temporal Cloud
func (w *workflows) GetAuditLogs(ctx workflow.Context, in *cloudservice.GetAuditLogsRequest) (*cloudservice.GetAuditLogsResponse, error) {
	return activities.GetAuditLogs(withInfiniteRetryActivityOptions(ctx), in) 
}
