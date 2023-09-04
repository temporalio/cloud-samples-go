package workflows

import (
	"fmt"
	"strings"
	"time"

	"github.com/temporalio/cloud-operations-workflows/internal/validator"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/operation/v1"
	"github.com/temporalio/cloud-operations-workflows/workflows/activities"
	"go.temporal.io/sdk/workflow"
)

type (
	WaitForAsyncOperationInput struct {
		AsyncOperationID string        `required:"true"`
		Timeout          time.Duration `required:"true"`
	}
	WaitForAsyncOperationOutput struct {
		AsyncOperation *operation.AsyncOperation
	}
)

// Get a async operation
func (w *workflows) GetAsynOperation(ctx workflow.Context, in *cloudservice.GetAsyncOperationRequest) (*cloudservice.GetAsyncOperationResponse, error) {
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
	getReqStatusFn := func(f workflow.Future) {
		resp, err = w.GetAsynOperation(ctx, &cloudservice.GetAsyncOperationRequest{
			AsyncOperationId: in.AsyncOperationID,
		})
	}

	// Check the request status immediately the first time, then poll at a regular interval afterwards
	selector.AddFuture(workflow.NewTimer(ctx, 0), getReqStatusFn)
	selector.AddFuture(workflow.NewTimer(ctx, in.Timeout), func(f workflow.Future) {
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
