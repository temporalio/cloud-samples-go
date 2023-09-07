package activities

import (
	"context"

	"github.com/gogo/status"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	CloudAPIRequestFailure = "temporal-cloud-api-request-failure"
)

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

func executeCloudAPIRequest[Req any, Resp any](
	ctx context.Context,
	in Req,
	fn func(context.Context, Req, ...grpc.CallOption) (Resp, error),
) (Resp, error) {
	out, err := fn(ctx, in)
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case
			codes.InvalidArgument,
			codes.NotFound,
			codes.AlreadyExists,
			codes.PermissionDenied,
			codes.ResourceExhausted,
			codes.FailedPrecondition,
			codes.Aborted,
			codes.OutOfRange,
			codes.Unimplemented,
			codes.Unauthenticated:

			// all these type of errors are application level errors and should fail the activity immediately
			return out, temporal.NewNonRetryableApplicationError(
				"CloudAPI request failed",
				CloudAPIRequestFailure,
				err,
			)
		}
	}
	// probably transient errors, let the activity retry
	return out, err
}
