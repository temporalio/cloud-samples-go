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

	CreatedReconcileResult     = "CREATED"
	UpdatedReconcileResult     = "UPDATED"
	DeletedReconcileResult     = "DELETED"
	UnchangedReconcileResult   = "UNCHANGED"
	UnaccountedReconcileResult = "UNACCOUNTED"
)

type (
	ReconcileUserResult struct {
		User   *user.User `json:"user"`
		Result string     `json:"result"`
	}

	ReconcileUserInput struct {
		Spec *user.UserSpec `required:"true" json:"spec"`
	}
	ReconcileUserOutput struct {
		Result ReconcileUserResult `json:"result"`
	}

	ReconcileUsersInput struct {
		Specs             []*user.UserSpec `required:"true" json:"specs"`
		DeleteUnaccounted bool             `json:"delete_unaccounted"`
	}
	ReconcileUsersOutput struct {
		Results []*ReconcileUserResult `json:"results"`
	}
	PeriodicReconcileUsersInput struct {
	}
	PeriodicReconcileUsersOutput struct {
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

func (w *workflows) reconcileUser(ctx workflow.Context, spec *user.UserSpec, user *user.User) (*user.User, string, error) {
	var (
		userID         string
		asyncOpID      string
		reconileStatus string
	)
	if user == nil {
		// no user found, create one
		createResp, err := w.CreateUser(ctx, &cloudservice.CreateUserRequest{
			Spec: spec,
		})
		if err != nil {
			return nil, "", err
		}
		userID = createResp.UserId
		asyncOpID = createResp.AsyncOperation.Id
		reconileStatus = CreatedReconcileResult

	} else if !proto.Equal(user.Spec, spec) {
		// user found, and specs don't match,  update it
		updateResp, err := w.UpdateUser(ctx, &cloudservice.UpdateUserRequest{
			UserId:          user.Id,
			Spec:            spec,
			ResourceVersion: user.ResourceVersion,
		})
		if err != nil {
			return nil, "", err
		}
		userID = user.Id
		asyncOpID = updateResp.AsyncOperation.Id
		reconileStatus = UpdatedReconcileResult
	} else {
		// nothing to change, get the latest user and return
		userID = user.Id
		reconileStatus = UnchangedReconcileResult
	}

	if asyncOpID != "" {
		// wait for the operation to complete
		_, err := w.WaitForAsyncOperation(ctx, &WaitForAsyncOperationInput{
			AsyncOperationID: asyncOpID,
			Timeout:          userUpdateTimeout,
		})
		if err != nil {
			return nil, "", err
		}
	}
	getResp, err := w.GetUser(ctx, &cloudservice.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, "", err
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
		Result: ReconcileUserResult{
			User:   user,
			Result: reconcileStatus,
		},
	}, nil
}

// Reconcile a user, create the user if one does not exist, or update the user if one does exist.
func (w *workflows) ReconcileUsers(ctx workflow.Context, in *ReconcileUsersInput) (*ReconcileUsersOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	var (
		getUsersReq = &cloudservice.GetUsersRequest{}
		users       = make(map[string]*user.User)
		out         = &ReconcileUsersOutput{}
	)
	for {
		resp, err := w.GetUsers(ctx, getUsersReq)
		if err != nil {
			return nil, err
		}
		for i := range resp.Users {
			users[resp.Users[i].Spec.Email] = resp.Users[i]
		}
		if resp.NextPageToken == "" {
			break
		}
		getUsersReq.PageToken = resp.NextPageToken
	}
	for i := range in.Specs {
		var user *user.User
		if u, ok := users[in.Specs[i].Email]; ok {
			user = u
		}
		// reconcile the user
		u, reconcileStatus, err := w.reconcileUser(ctx, in.Specs[i], user)
		if err != nil {
			return nil, err
		}
		out.Results = append(out.Results, &ReconcileUserResult{
			User:   u,
			Result: reconcileStatus,
		})
		// remove the reconciled users from the map
		delete(users, in.Specs[i].Email)
	}

	// whats left in maps is only the unaccounted users
	for _, u := range users {
		if in.DeleteUnaccounted {
			_, err := w.DeleteUser(ctx, &cloudservice.DeleteUserRequest{
				UserId:          u.Id,
				ResourceVersion: u.ResourceVersion,
			})
			if err != nil {
				return nil, err
			}
			out.Results = append(out.Results, &ReconcileUserResult{
				User:   u,
				Result: DeletedReconcileResult,
			})
		} else {
			out.Results = append(out.Results, &ReconcileUserResult{
				User:   u,
				Result: UnaccountedReconcileResult,
			})
		}
	}
	return out, nil
}
