package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	defaultActivityStartToCloseTimeout = 1 * time.Minute
)

func withActivityStartToCloseTimeout(timeout time.Duration) func(*workflow.ActivityOptions) {
	return func(ao *workflow.ActivityOptions) {
		ao.StartToCloseTimeout = timeout
	}
}

func withActivityScheduleToCloseTimeout(timeout time.Duration) func(*workflow.ActivityOptions) {
	return func(ao *workflow.ActivityOptions) {
		ao.ScheduleToCloseTimeout = timeout
	}
}

func withActivityHeartbeatTimeout(timeout time.Duration) func(*workflow.ActivityOptions) {
	return func(ao *workflow.ActivityOptions) {
		ao.HeartbeatTimeout = timeout
	}
}

func withActivityRetryPolicy(retryPolicy *temporal.RetryPolicy) func(*workflow.ActivityOptions) {
	return func(ao *workflow.ActivityOptions) {
		ao.RetryPolicy = retryPolicy
	}
}

// withInfiniteRetryActivityOptions returns a context with activity options allowing infinite retry with each attempt
// subject to defaultActivityStartToCloseTimeout.
//
// - NOTE: This is the preferred default setup for activity options
// - Use withActivityStartToCloseTimeout to override defaultActivityStartToCloseTimeout
// - Use withActivityHeartbeatTimeout to pass in a heartbeat timeout
func withInfiniteRetryActivityOptions(ctx workflow.Context, options ...func(*workflow.ActivityOptions)) workflow.Context {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultActivityStartToCloseTimeout,
	}
	for _, option := range options {
		option(&ao)
	}
	return workflow.WithActivityOptions(ctx, ao)
}
