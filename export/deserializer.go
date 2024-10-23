package export

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/export/v1"
)

// DeserializeExportedWorkflows deserializes a byte array into a WorkflowExecutions object. This is useful for programmatically processing workflow histories
func DeserializeExportedWorkflows(bytes []byte) (*export.WorkflowExecutions, error) {
	blob := &common.DataBlob{
		EncodingType: enumspb.ENCODING_TYPE_PROTO3, // Currently workflow histories are only exported in proto3 format
		Data:         bytes,
	}

	var workflows export.WorkflowExecutions
	if err := proto.Unmarshal(blob.Data, &workflows); err != nil {
		return nil, fmt.Errorf("failed to decode export file: %w", err)
	}

	return &workflows, nil
}

// FormatWorkflow converts an exported workflow execution into a friendly, human-readable string
func FormatWorkflow(workflow *export.WorkflowExecution) string {
	pbMarshaler := protojson.MarshalOptions{
		Indent:            "\t",
		EmitDefaultValues: true,
	}
	return pbMarshaler.Format(workflow)
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
