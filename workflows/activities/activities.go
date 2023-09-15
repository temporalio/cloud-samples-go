package activities

import (
	"github.com/temporalio/cloud-samples-go/protogen/temporal/api/cloud/cloudservice/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	activitiesPrefix = "temporal-cloud-activities."
)

type (
	Activities struct {
		cloudserviceclient cloudservice.CloudServiceClient
	}
)

func NewActivities(conn grpc.ClientConnInterface) *Activities {
	return &Activities{cloudserviceclient: cloudservice.NewCloudServiceClient(conn)}
}

func Register(w worker.Worker, activities *Activities) {
	w.RegisterActivityWithOptions(activities, activity.RegisterOptions{Name: activitiesPrefix})
}
