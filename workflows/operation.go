package workflows

import (
	"fmt"
	"strings"
	"time"

	"github.com/temporalio/cloud-samples-go/internal/validator"
	"go.temporal.io/api/cloud/cloudservice/v1"
	"go.temporal.io/api/cloud/operation/v1"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	// async operation workflow types
	GetAsyncOperationWorkflowType = workflowPrefix + "get-async-operation"
	WaitForAsyncOperationType     = workflowPrefix + "wait-for-async-operation"
)

type (
	WaitForAsyncOperationInput struct {
		AsyncOperationID string        `required:"true"`
		Timeout          time.Duration `required:"true"`
	}
	WaitForAsyncOperationOutput struct {
		AsyncOperation *operation.AsyncOperation
	}

	AsyncOperationWorkflows interface {
		// Async Operations Workflows
		GetAsyncOperation(ctx workflow.Context, in *cloudservice.GetAsyncOperationRequest) (*cloudservice.GetAsyncOperationResponse, error)
		WaitForAsyncOperation(ctx workflow.Context, in *WaitForAsyncOperationInput) (*WaitForAsyncOperationOutput, error)
	}
)

func registerAsyncOperationWorkflows(w worker.Worker, wf AsyncOperationWorkflows) {
	for k, v := range map[string]any{
		GetAsyncOperationWorkflowType: wf.GetAsyncOperation,
		WaitForAsyncOperationType:     wf.WaitForAsyncOperation,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}
}

// Get a async operation
func (w *workflows) GetAsyncOperation(ctx workflow.Context, in *cloudservice.GetAsyncOperationRequest) (*cloudservice.GetAsyncOperationResponse, error) {
	return activities.GetAsyncOperation(withInfiniteRetryActivityOptions(ctx), in)
}

// Wait for the async operation to finish
func (w *workflows) WaitForAsyncOperation(ctx workflow.Context, in *WaitForAsyncOperationInput) (*WaitForAsyncOperationOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	var (
		resp *cloudservice.GetAsyncOperationResponse
		err  error
	)
	selector := workflow.NewSelector(ctx)
	getReqStatusFn := func(_ workflow.Future) {
		resp, err = w.GetAsyncOperation(ctx, &cloudservice.GetAsyncOperationRequest{
			AsyncOperationId: in.AsyncOperationID,
		})
	}

	// Check the request status immediately the first time, then poll at a regular interval afterwards
	selector.AddFuture(workflow.NewTimer(ctx, 0), getReqStatusFn)
	selector.AddFuture(workflow.NewTimer(ctx, in.Timeout), func(_ workflow.Future) {
		err = fmt.Errorf("timed out waiting for async operation, asyncOperationID=%s, timeout=%s",
			in.AsyncOperationID, in.Timeout)
	})
	for {
		selector.Select(ctx)
		if err != nil {
			return nil, err
		}
		switch {
		case strings.EqualFold(resp.AsyncOperation.State, "FAILED"):
			return nil, fmt.Errorf("request failed: %s", resp.AsyncOperation.FailureReason)
		case strings.EqualFold(resp.AsyncOperation.State, "FULFILLED"):
			return &WaitForAsyncOperationOutput{
				AsyncOperation: resp.AsyncOperation,
			}, nil
		default:
			selector.AddFuture(workflow.NewTimer(ctx, resp.AsyncOperation.CheckDuration.AsDuration()), getReqStatusFn)
		}
	}
}
