package activities

import (
	"context"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
)

func (a *Activities) GetAuditLogs(ctx context.Context, in *cloudservice.GetAuditLogsRequest) (*cloudservice.GetAuditLogsResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetAuditLogs)
}

var (
	GetAuditLogs = executeActivityFn[*cloudservice.GetAuditLogsRequest, *cloudservice.GetAuditLogsResponse](activitiesPrefix + "GetAuditLogs")
)
