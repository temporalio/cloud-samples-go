package activities

import (
	"github.com/temporalio/cloud-samples-go/client/api"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

const (
	activitiesPrefix = "tmprlcloud-activity."
)

type (
	Activities struct {
		client *api.Client
	}
)

func NewActivities(client *api.Client) *Activities {
	return &Activities{client: client}
}

func Register(w worker.Worker, activities *Activities) {
	w.RegisterActivityWithOptions(activities, activity.RegisterOptions{Name: activitiesPrefix})
}
