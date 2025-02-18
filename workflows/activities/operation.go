package activities

import (
	"context"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
)

func (a *Activities) GetAsyncOperation(ctx context.Context, in *cloudservice.GetAsyncOperationRequest) (*cloudservice.GetAsyncOperationResponse, error) {
	return a.client.CloudService().GetAsyncOperation(ctx, in)
}

var GetAsyncOperation = executeActivityFn[*cloudservice.GetAsyncOperationRequest, *cloudservice.GetAsyncOperationResponse](activitiesPrefix + "GetAsyncOperation")
