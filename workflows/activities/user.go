package activities

import (
	"context"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
)

func (a *Activities) GetUser(ctx context.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetUser)
}

func (a *Activities) GetUsers(ctx context.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().GetUsers)
}

func (a *Activities) CreateUser(ctx context.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().CreateUser)
}

func (a *Activities) UpdateUser(ctx context.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().UpdateUser)
}

func (a *Activities) DeleteUser(ctx context.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.client.CloudService().DeleteUser)
}

var (
	GetUser    = executeActivityFn[*cloudservice.GetUserRequest, *cloudservice.GetUserResponse](activitiesPrefix + "GetUser")
	GetUsers   = executeActivityFn[*cloudservice.GetUsersRequest, *cloudservice.GetUsersResponse](activitiesPrefix + "GetUsers")
	CreateUser = executeActivityFn[*cloudservice.CreateUserRequest, *cloudservice.CreateUserResponse](activitiesPrefix + "CreateUser")
	UpdateUser = executeActivityFn[*cloudservice.UpdateUserRequest, *cloudservice.UpdateUserResponse](activitiesPrefix + "UpdateUser")
	DeleteUser = executeActivityFn[*cloudservice.DeleteUserRequest, *cloudservice.DeleteUserResponse](activitiesPrefix + "DeleteUser")
)
