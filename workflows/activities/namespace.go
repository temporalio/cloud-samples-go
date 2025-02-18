package activities

import (
	"context"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
)

func (a *Activities) GetNamespace(ctx context.Context, in *cloudservice.GetNamespaceRequest) (*cloudservice.GetNamespaceResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetNamespace)
}

func (a *Activities) GetNamespaces(ctx context.Context, in *cloudservice.GetNamespacesRequest) (*cloudservice.GetNamespacesResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetNamespaces)
}

func (a *Activities) CreateNamespace(ctx context.Context, in *cloudservice.CreateNamespaceRequest) (*cloudservice.CreateNamespaceResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().CreateNamespace)
}

func (a *Activities) UpdateNamespace(ctx context.Context, in *cloudservice.UpdateNamespaceRequest) (*cloudservice.UpdateNamespaceResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().UpdateNamespace)
}

func (a *Activities) DeleteNamespace(ctx context.Context, in *cloudservice.DeleteNamespaceRequest) (*cloudservice.DeleteNamespaceResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().DeleteNamespace)
}

var (
	GetNamespace    = executeActivityFn[*cloudservice.GetNamespaceRequest, *cloudservice.GetNamespaceResponse](activitiesPrefix + "GetNamespace")
	GetNamespaces   = executeActivityFn[*cloudservice.GetNamespacesRequest, *cloudservice.GetNamespacesResponse](activitiesPrefix + "GetNamespaces")
	CreateNamespace = executeActivityFn[*cloudservice.CreateNamespaceRequest, *cloudservice.CreateNamespaceResponse](activitiesPrefix + "CreateNamespace")
	UpdateNamespace = executeActivityFn[*cloudservice.UpdateNamespaceRequest, *cloudservice.UpdateNamespaceResponse](activitiesPrefix + "UpdateNamespace")
	DeleteNamespace = executeActivityFn[*cloudservice.DeleteNamespaceRequest, *cloudservice.DeleteNamespaceResponse](activitiesPrefix + "DeleteNamespace")
)
