package activities

import "go.temporal.io/sdk/workflow"

type (
	ExecuteActivity[Req any, Resp any] func(workflow.Context, Req) (Resp, error)
)

func executeActivityFn[Req any, Resp any](activityName string) ExecuteActivity[Req, Resp] {
	return func(ctx workflow.Context, in Req) (Resp, error) {
		var out Resp
		err := workflow.ExecuteActivity(ctx, activityName, in).Get(ctx, &out)
		return out, err
	}
}
