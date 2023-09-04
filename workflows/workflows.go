package workflows

import (
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/cloudservice/v1"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/user/v1"
	"github.com/temporalio/cloud-operations-workflows/workflows/activities"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
)

//go:generate mockgen -source workflows.go -destination workflows_mock.go -package workflow

const (
	workflowPrefix         = "cloud-operations-workflows."
	GetUserWorkflowType    = workflowPrefix + "get-user"
	GetUsersWorkflowType   = workflowPrefix + "get-users"
	CreateUserWorkflowType = workflowPrefix + "create-user"
	UpdateUserWorkflowType = workflowPrefix + "update-user"
	DeleteUserWorkflowType = workflowPrefix + "delete-user"
)

type (
	Workflows interface {
		// User Management
		GetUser(ctx workflow.Context, in *cloudservice.GetUserRequest) (*cloudservice.GetUserResponse, error)
		GetUsers(ctx workflow.Context, in *cloudservice.GetUsersRequest) (*cloudservice.GetUsersResponse, error)
		GetUserWithEmail(ctx workflow.Context, email string) (*user.User, error)
		CreateUser(ctx workflow.Context, in *cloudservice.CreateUserRequest) (*cloudservice.CreateUserResponse, error)
		UpdateUser(ctx workflow.Context, in *cloudservice.UpdateUserRequest) (*cloudservice.UpdateUserResponse, error)
		DeleteUser(ctx workflow.Context, in *cloudservice.DeleteUserRequest) (*cloudservice.DeleteUserResponse, error)
		ReconcileUser(ctx workflow.Context, in *ReconcileUserInput) (*ReconcileUserOutput, error)
	}

	workflows struct{}
)

func NewWorkflows() Workflows {
	return &workflows{}
}

func NewActivities(conn grpc.ClientConnInterface) *activities.Activities {
	return activities.NewActivities(conn)
}

func Register(w worker.Worker, wf Workflows, a *activities.Activities) {
	// Register the workflows that we want to be able to use.
	for k, v := range map[string]any{
		GetUserWorkflowType:    wf.GetUser,
		GetUsersWorkflowType:   wf.GetUsers,
		CreateUserWorkflowType: wf.CreateUser,
		UpdateUserWorkflowType: wf.UpdateUser,
		DeleteUserWorkflowType: wf.DeleteUser,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}

	// Register the activities that the workflows will use.
	activities.Register(w, a)
}
