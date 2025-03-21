package activities

import (
	"context"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
)

func (a *Activities) GetRegion(ctx context.Context, in *cloudservice.GetRegionRequest) (*cloudservice.GetRegionResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetRegion)
}

func (a *Activities) GetRegions(ctx context.Context, in *cloudservice.GetRegionsRequest) (*cloudservice.GetRegionsResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetRegions)
}

var (
	GetRegion  = executeActivityFn[*cloudservice.GetRegionRequest, *cloudservice.GetRegionResponse](activitiesPrefix + "GetRegion")
	GetRegions = executeActivityFn[*cloudservice.GetRegionsRequest, *cloudservice.GetRegionsResponse](activitiesPrefix + "GetRegions")
)
