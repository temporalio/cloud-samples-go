package temporal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/cloud-samples-go/client/temporal"
	"github.com/temporalio/cloud-samples-go/protogen/temporal/api/workflowservice/v1"
)

func TestCloudClient(t *testing.T) {

	assert := assert.New(t)
	client, err := temporal.GetTemporalCloudNamespaceClient(&temporal.GetTemporalCloudNamespaceClientInput{
		Namespace:       "abhinav-test3.a2dd6",
		TLSCertFilePath: "/Users/abhinavtemporal.io/development/certs/test/cert.pem",
		TLSKeyFilePath:  "/Users/abhinavtemporal.io/development/certs/test/cert.key",
	})
	assert.NoError(err)
	defer client.Close()

	_, err = client.ListOpenWorkflow(context.Background(), &workflowservice.ListOpenWorkflowExecutionsRequest{})
	assert.NoError(err)
}
