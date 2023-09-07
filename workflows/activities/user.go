package activities

import (
	"context"

	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
)

func (a *Activities) GetUser(ctx context.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.GetUser)
}

func (a *Activities) GetUsers(ctx context.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.GetUsers)
}

func (a *Activities) CreateUser(ctx context.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.CreateUser)
}

func (a *Activities) UpdateUser(ctx context.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.UpdateUser)
}

func (a *Activities) DeleteUser(ctx context.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error) {
	return executeCloudAPIRequest(ctx, in, a.cloudserviceclient.DeleteUser)
}

var (
	GetUser    = executeActivityFn[*cloudservice.GetUserRequest, *cloudservice.GetUserResponse](activitiesPrefix + "GetUser")
	GetUsers   = executeActivityFn[*cloudservice.GetUsersRequest, *cloudservice.GetUsersResponse](activitiesPrefix + "GetUsers")
	CreateUser = executeActivityFn[*cloudservice.CreateUserRequest, *cloudservice.CreateUserResponse](activitiesPrefix + "CreateUser")
	UpdateUser = executeActivityFn[*cloudservice.UpdateUserRequest, *cloudservice.UpdateUserResponse](activitiesPrefix + "UpdateUser")
	DeleteUser = executeActivityFn[*cloudservice.DeleteUserRequest, *cloudservice.DeleteUserResponse](activitiesPrefix + "DeleteUser")
)
