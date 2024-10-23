package activities

import (
	"context"

	"go.temporal.io/api/cloud/cloudservice/v1"
)

func (a *Activities) GetRegion(ctx context.Context, in *cloudservice.GetRegionRequest) (*cloudservice.GetRegionResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.GetRegion)
}

func (a *Activities) GetRegions(ctx context.Context, in *cloudservice.GetRegionsRequest) (*cloudservice.GetRegionsResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.GetRegions)
}

var (
	GetRegion  = executeActivityFn[*cloudservice.GetRegionRequest, *cloudservice.GetRegionResponse](activitiesPrefix + "GetRegion")
	GetRegions = executeActivityFn[*cloudservice.GetRegionsRequest, *cloudservice.GetRegionsResponse](activitiesPrefix + "GetRegions")
)
