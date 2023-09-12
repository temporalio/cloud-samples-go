package activities

import (
	"context"
	"encoding/json"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/user/v1"
	"go.temporal.io/sdk/activity"
	"google.golang.org/grpc/status"
	"io"
	"os"
)

type (
	GetUserSpecsFromFileRequest struct {
		FilePath string
	}
	GetUserSpecsFromFileResponse struct {
		Specs []*user.UserSpec `json:"users"`
	}
)

func (a *Activities) GetUser(ctx context.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error) {
	return a.cloudserviceclient.GetUser(ctx, in)
}

func (a *Activities) GetUsers(ctx context.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error) {
	return a.cloudserviceclient.GetUsers(ctx, in)
}

func (a *Activities) CreateUser(ctx context.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error) {
	return a.cloudserviceclient.CreateUser(ctx, in)
}

func (a *Activities) UpdateUser(ctx context.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error) {
	resp, err := a.cloudserviceclient.UpdateUser(ctx, in)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Message() == "nothing to change" {
				return nil, nil
			}
		}
		return nil, err
	}
	return resp, nil
}

func (a *Activities) DeleteUser(ctx context.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error) {
	return a.cloudserviceclient.DeleteUser(ctx, in)
}

func (a *Activities) GetUserSpecsFromFile(ctx context.Context, in *GetUserSpecsFromFileRequest) (*GetUserSpecsFromFileResponse, error) {
	jsonFile, err := os.Open(in.FilePath)

	if err != nil {
		activity.GetLogger(ctx).Error("Error opening file", "error", err)
		return nil, err
	}
	//activity.GetLogger(ctx).Info("Successfully Opened file " + in.FilePath)
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		activity.GetLogger(ctx).Error("Error reading file", "error", err)
		return nil, err
	}

	var resp GetUserSpecsFromFileResponse
	err = json.Unmarshal(byteValue, &resp)
	return &resp, err
}

var (
	GetUser      = executeActivityFn[*cloudservice.GetUserRequest, *cloudservice.GetUserResponse](activitiesPrefix + "GetUser")
	GetUsers     = executeActivityFn[*cloudservice.GetUsersRequest, *cloudservice.GetUsersResponse](activitiesPrefix + "GetUsers")
	CreateUser   = executeActivityFn[*cloudservice.CreateUserRequest, *cloudservice.CreateUserResponse](activitiesPrefix + "CreateUser")
	UpdateUser   = executeActivityFn[*cloudservice.UpdateUserRequest, *cloudservice.UpdateUserResponse](activitiesPrefix + "UpdateUser")
	DeleteUser   = executeActivityFn[*cloudservice.DeleteUserRequest, *cloudservice.DeleteUserResponse](activitiesPrefix + "DeleteUser")
	GetUserSpecs = executeActivityFn[*GetUserSpecsFromFileRequest, *GetUserSpecsFromFileResponse](activitiesPrefix + "GetUserSpecsFromFile")
)
