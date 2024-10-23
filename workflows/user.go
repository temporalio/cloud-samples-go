package workflows

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/proto"

	"github.com/temporalio/cloud-samples-go/internal/validator"
	"go.temporal.io/api/cloud/cloudservice/v1"
	"go.temporal.io/api/cloud/identity/v1"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
)

const (
	userUpdateTimeout = 10 * time.Minute

	// user management workflow types
	GetUserWorkflowType                      = workflowPrefix + "get-user"
	GetUsersWorkflowType                     = workflowPrefix + "get-users"
	GetAllUsersWorkflowType                  = workflowPrefix + "get-all-users"
	GetUserWithEmailWorkflow                 = workflowPrefix + "get-user-with-email"
	GetAllUsersWithAccessToNamespaceWorkflow = workflowPrefix + "get-all-users-with-access-to-namespace"
	CreateUserWorkflowType                   = workflowPrefix + "create-user"
	UpdateUserWorkflowType                   = workflowPrefix + "update-user"
	DeleteUserWorkflowType                   = workflowPrefix + "delete-user"
	ReconcileUserWorkflowType                = workflowPrefix + "reconcile-user"
	ReconcileUsersWorkflowType               = workflowPrefix + "reconcile-users"

	// reconcile outcomes
	ReconcileOutcomeCreated   = "created"
	ReconcileOutcomeDeleted   = "deleted"
	ReconcileOutcomeUpdated   = "updated"
	ReconcileOutcomeUnchanged = "unchanged"
	ReconcileOutcomeError     = "error"
)

type (
	ReconcileOutcome   string
	ReconcileUserInput struct {
		Spec *identity.UserSpec `required:"true" json:"spec"`
	}
	ReconcileUserOutput struct {
		User    *identity.User   `json:"user"`
		Outcome ReconcileOutcome `json:"outcome"`
		Error   string           `json:"error"`
	}

	ReconcileUsersInput struct {
		Specs             []*identity.UserSpec `required:"true" json:"specs"`
		DeleteUnaccounted bool                 `json:"delete_unaccounted"`
	}
	ReconcileUsersOutput struct {
		Results []*ReconcileUserOutput `json:"results"`
	}

	UserWorkflows interface {
		// User Management Workflows
		GetUser(ctx workflow.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error)
		GetUsers(ctx workflow.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error)
		GetAllUsers(ctx workflow.Context) ([]*identity.User, error)
		GetUserWithEmail(ctx workflow.Context, email string) (*identity.User, error)
		GetAllUsersWithAccessToNamespace(ctx workflow.Context, namespace string) ([]*identity.User, error)
		CreateUser(ctx workflow.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error)
		UpdateUser(ctx workflow.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error)
		DeleteUser(ctx workflow.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error)
		ReconcileUser(ctx workflow.Context, in *ReconcileUserInput) (*ReconcileUserOutput, error)
		ReconcileUsers(ctx workflow.Context, in *ReconcileUsersInput) (*ReconcileUsersOutput, error)
	}
)

func (o *ReconcileUserOutput) setError(err error) {
	var applicationErr *temporal.ApplicationError
	if errors.As(err, &applicationErr) {
		o.Error = applicationErr.Error()
		o.Outcome = ReconcileOutcomeError
	}
}

func registerUserWorkflows(w worker.Worker, wf UserWorkflows) {
	for k, v := range map[string]any{
		GetUserWorkflowType:                      wf.GetUser,
		GetUsersWorkflowType:                     wf.GetUsers,
		GetAllUsersWorkflowType:                  wf.GetAllUsers,
		GetUserWithEmailWorkflow:                 wf.GetUserWithEmail,
		GetAllUsersWithAccessToNamespaceWorkflow: wf.GetAllUsersWithAccessToNamespace,
		CreateUserWorkflowType:                   wf.CreateUser,
		UpdateUserWorkflowType:                   wf.UpdateUser,
		DeleteUserWorkflowType:                   wf.DeleteUser,
		ReconcileUserWorkflowType:                wf.ReconcileUser,
		ReconcileUsersWorkflowType:               wf.ReconcileUsers,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}
}

// Get a user
func (w *workflows) GetUser(ctx workflow.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error) {
	return activities.GetUser(withInfiniteRetryActivityOptions(ctx), in)
}

// Get multiple users
func (w *workflows) GetUsers(ctx workflow.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error) {
	return activities.GetUsers(withInfiniteRetryActivityOptions(ctx), in)
}

func (w *workflows) getAllUsers(ctx workflow.Context, email, namespace string) ([]*identity.User, error) {
	var (
		users     = make([]*identity.User, 0)
		pageToken = ""
	)
	for {
		resp, err := w.GetUsers(ctx, &cloudservice.GetUsersRequest{
			Email:     email,
			Namespace: namespace,
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}
		users = append(users, resp.Users...)
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}
	return users, nil
}

// Get all known Users
func (w *workflows) GetAllUsers(ctx workflow.Context) ([]*identity.User, error) {
	return w.getAllUsers(ctx, "", "")
}

// Get the user with email
func (w *workflows) GetUserWithEmail(ctx workflow.Context, email string) (*identity.User, error) {
	users, err := w.getAllUsers(ctx, email, "")
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("multiple users found for email %s", email)
	}
	return users[0], nil
}

// Get all the users who have access to namespace
func (w *workflows) GetAllUsersWithAccessToNamespace(ctx workflow.Context, namespace string) ([]*identity.User, error) {
	return w.getAllUsers(ctx, "", namespace)
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

func (w *workflows) reconcileUser(ctx workflow.Context, spec *identity.UserSpec, user *identity.User) (*ReconcileUserOutput, error) {
	var (
		userID    string
		asyncOpID string
		out       = &ReconcileUserOutput{}
		err       error
	)
	defer func() {
		if err != nil {
			out.setError(err)
		}
		if user != nil {
			out.User = user
		} else if spec != nil {
			out.User = &identity.User{
				Id:   userID,
				Spec: spec,
			}
		}
	}()
	if user == nil {
		var createResp *cloudservice.CreateUserResponse
		// no user found, create one
		createResp, err = w.CreateUser(ctx, &cloudservice.CreateUserRequest{
			Spec: spec,
		})
		if err != nil {
			return out, err
		}
		userID = createResp.UserId
		asyncOpID = createResp.AsyncOperation.Id
		out.Outcome = ReconcileOutcomeCreated

	} else if !proto.Equal(user.Spec, spec) {
		var updateResp *cloudservice.UpdateUserResponse
		// user found, and specs don't match,  update it
		updateResp, err = w.UpdateUser(ctx, &cloudservice.UpdateUserRequest{
			UserId:          user.Id,
			Spec:            spec,
			ResourceVersion: user.ResourceVersion,
		})
		if err != nil {
			return out, err
		}
		userID = user.Id
		asyncOpID = updateResp.AsyncOperation.Id
		out.Outcome = ReconcileOutcomeUpdated

	} else {
		// nothing to change, get the latest user and return
		userID = user.Id
		out.Outcome = ReconcileOutcomeUnchanged
		return out, nil
	}

	if asyncOpID != "" {
		// wait for the operation to complete
		_, err = w.WaitForAsyncOperation(ctx, &WaitForAsyncOperationInput{
			AsyncOperationID: asyncOpID,
			Timeout:          userUpdateTimeout,
		})
		if err != nil {
			return out, err
		}
	}
	var getResp *cloudservice.GetUserResponse
	getResp, err = w.GetUser(ctx, &cloudservice.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return out, err
	}
	user = getResp.User
	return out, nil
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
	out, err := w.reconcileUser(ctx, in.Spec, user)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Reconcile multiple users, create missing users, update existing users, and optionally delete unaccounted users.
func (w *workflows) ReconcileUsers(ctx workflow.Context, in *ReconcileUsersInput) (*ReconcileUsersOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	users, err := w.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	out := &ReconcileUsersOutput{}
	for i := range in.Specs {
		var user *identity.User
		for _, u := range users {
			if u.Spec.Email == in.Specs[i].Email {
				user = u
				users = append(users[:i], users[i+1:]...) // remove the user from the list
				break
			}
		}
		// reconcile the user
		o, _ := w.reconcileUser(ctx, in.Specs[i], user)
		out.Results = append(out.Results, o)
	}
	// whats left in maps is only the unaccounted users
	for _, u := range users {
		if in.DeleteUnaccounted {
			o := &ReconcileUserOutput{
				User:    u,
				Outcome: ReconcileOutcomeDeleted,
			}
			_, err := w.DeleteUser(ctx, &cloudservice.DeleteUserRequest{
				UserId:          u.Id,
				ResourceVersion: u.ResourceVersion,
			})
			if err != nil {
				o.setError(err)
			}
			out.Results = append(out.Results, o)
		}
	}
	return out, nil
}
