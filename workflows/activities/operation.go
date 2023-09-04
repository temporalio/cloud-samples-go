package activities

import (
	"context"

	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
)

func (a *Activities) GetAsyncOperation(ctx context.Context, in *cloudservice.GetAsyncOperationRequest) (*cloudservice.GetAsyncOperationResponse, error) {
	return a.cloudserviceclient.GetAsyncOperation(ctx, in)
}

var (
	GetAsyncOperation = executeActivityFn[*cloudservice.GetAsyncOperationRequest, *cloudservice.GetAsyncOperationResponse](activitiesPrefix + "GetAsyncOperation")
)
