package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/proto"

	"github.com/temporalio/cloud-operations-workflows/internal/validator"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/user/v1"
	"github.com/temporalio/cloud-operations-workflows/workflows/activities"
)

const (
	userUpdateTimeout = 10 * time.Minute

	CreatedReconcileUserStatus   = ReconcileUserStatus("CREATED")
	UpdatedReconcileUserStatus   = ReconcileUserStatus("UPDATED")
	DeletedReconcileUserStatus   = ReconcileUserStatus("DELETED")
	UnchangedReconcileUserStatus = ReconcileUserStatus("UNCHANGED")
)

type (
	ReconcileUserStatus string

	ReconcileUserInput struct {
		Spec *user.UserSpec `required:"true"`
	}
	ReconcileUserOutput struct {
		User   *user.User
		Status ReconcileUserStatus
	}
)

// Get a user
func (w *workflows) GetUser(ctx workflow.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error) {
	return activities.GetUser(withInfiniteRetryActivityOptions(ctx), in)
}

// Get multiple users
func (w *workflows) GetUsers(ctx workflow.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error) {
	return activities.GetUsers(withInfiniteRetryActivityOptions(ctx), in)
}

// Get the user with email
func (w *workflows) GetUserWithEmail(ctx workflow.Context, email string) (*user.User, error) {
	resp, err := w.GetUsers(ctx, &cloudservice.GetUsersRequest{
		EmailAddress: email,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, nil
	}
	if len(resp.Users) > 1 {
		return nil, fmt.Errorf("multiple users found for email %s", email)
	}
	return resp.Users[0], nil
}

// Create a user
func (w *workflows) CreateUser(ctx workflow.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error) {
	return activities.CreateUser(withInfiniteRetryActivityOptions(ctx), in)
}

// Update a user
func (w *workflows) UpdateUser(ctx workflow.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error) {
	return activities.UpdateUser(withInfiniteRetryActivityOptions(ctx), in)
}

// Delete a user
func (w *workflows) DeleteUser(ctx workflow.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error) {
	return activities.DeleteUser(withInfiniteRetryActivityOptions(ctx), in)
}

func (w *workflows) reconcileUser(ctx workflow.Context, spec *user.UserSpec, user *user.User) (*user.User, ReconcileUserStatus, error) {
	var (
		userID         string
		asyncOpID      string
		reconileStatus = UnchangedReconcileUserStatus
	)
	if user == nil {
		// no user found, create one
		createResp, err := w.CreateUser(ctx, &cloudservice.CreateUserRequest{})
		if err != nil {
			return nil, reconileStatus, err
		}
		userID = createResp.UserId
		asyncOpID = createResp.AsyncOperation.Id
		reconileStatus = CreatedReconcileUserStatus

	} else if !proto.Equal(user.Spec, spec) {
		// user found, and specs don't match,  update it
		updateResp, err := w.UpdateUser(ctx, &cloudservice.UpdateUserRequest{
			UserId:          user.Id,
			Spec:            spec,
			ResourceVersion: user.ResourceVersion,
		})
		if err != nil {
			return nil, reconileStatus, err
		}
		userID = user.Id
		asyncOpID = updateResp.AsyncOperation.Id
		reconileStatus = UpdatedReconcileUserStatus
	}

	// wait for the operation to complete
	_, err := w.WaitForAsyncOperation(ctx, &WaitForAsyncOperationInput{
		AsyncOperationID: asyncOpID,
		Timeout:          userUpdateTimeout,
	})
	if err != nil {
		return nil, reconileStatus, err
	}
	getResp, err := w.GetUser(ctx, &cloudservice.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, reconileStatus, err
	}
	return getResp.User, reconileStatus, nil
}

// Reconcile a user, create the user if one does not exist, or update the user if one does exist.
func (w *workflows) ReconcileUser(ctx workflow.Context, in *ReconcileUserInput) (*ReconcileUserOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	user, err := w.GetUserWithEmail(ctx, in.Spec.Email)
	if err != nil {
		return nil, err
	}
	user, reconcileStatus, err := w.reconcileUser(ctx, in.Spec, user)
	if err != nil {
		return nil, err
	}
	return &ReconcileUserOutput{
		User:   user,
		Status: reconcileStatus,
	}, nil
}
