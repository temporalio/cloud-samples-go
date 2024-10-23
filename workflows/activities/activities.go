package activities

import (
	"go.temporal.io/api/cloud/cloudservice/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

const (
	activitiesPrefix = "tmprlcloud-activity."
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
