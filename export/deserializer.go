package export

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"go.temporal.io/api/export/v1"
)

// DeserializeExportedExecutions decodes a Protobuf export file containing workflow and activity executions.
func DeserializeExportedExecutions(bytes []byte) (*export.ExportedExecutions, error) {
	var executions export.ExportedExecutions
	if err := proto.Unmarshal(bytes, &executions); err != nil {
		return nil, fmt.Errorf("failed to decode export file: %w", err)
	}

	return &executions, nil
}

// FormatWorkflow converts an exported workflow execution into a friendly, human-readable string
func FormatWorkflow(workflow *export.WorkflowExecution) string {
	pbMarshaler := protojson.MarshalOptions{
		Indent:            "\t",
		EmitDefaultValues: true,
	}
	return pbMarshaler.Format(workflow)
}

// FormatActivity converts an exported activity execution into a friendly, human-readable string.
func FormatActivity(activity *export.ActivityExecution) string {
	pbMarshaler := protojson.MarshalOptions{
		Indent:            "\t",
		EmitDefaultValues: true,
	}
	return pbMarshaler.Format(activity)
}

// GetExportedActivityInformation returns a string containing the activity ID, run ID, and activity type.
func GetExportedActivityInformation(activity *export.ActivityExecution) (string, error) {
	execution := activity.GetActivity()
	if execution == nil {
		return "", fmt.Errorf("activity execution is nil")
	}

	info := execution.GetInfo()
	if info == nil {
		return "", fmt.Errorf("activity execution information is nil")
	}

	return fmt.Sprintf("ActivityID: %s, RunID: %s, ActivityType: %s", info.GetActivityId(), info.GetRunId(), info.GetActivityType().GetName()), nil
}

// GetExportedWorkflowInformation returns a string containing the workflow ID, run ID, and workflow type
func GetExportedWorkflowInformation(workflow *export.WorkflowExecution) (string, error) {
	history := workflow.GetHistory()
	if history == nil {
		return "", fmt.Errorf("workflow history is nil")
	}

	events := history.GetEvents()
	if len(events) == 0 {
		return "", fmt.Errorf("workflow history has no events")
	}

	firstEvent := events[0]
	startAttributes := firstEvent.GetWorkflowExecutionStartedEventAttributes()
	if startAttributes == nil {
		return "", fmt.Errorf("first workflow history is not a start event")
	}

	return fmt.Sprintf("WorkflowID: %s, RunID: %s, WorkflowType: %s", startAttributes.GetWorkflowId(), startAttributes.GetOriginalExecutionRunId(), startAttributes.GetWorkflowType().GetName()), nil
}
