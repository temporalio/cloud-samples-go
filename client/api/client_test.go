package api

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/cloud-samples-go/protogen/temporal/api/cloud/cloudservice/v1"
)

const (
	temporalCloudAPIAddress    = "saas-api.tmprl.cloud:443"
	temporalCloudAPIKeyEnvName = "TEMPORAL_CLOUD_API_KEY"
)

func getAPIKeyFromEnv() (string, error) {
	v := os.Getenv(temporalCloudAPIKeyEnvName)
	if v == "" {
		return "", fmt.Errorf("apikey not provided, set environment variable '%s' with apikey you want to use", temporalCloudAPIKeyEnvName)
	}
	return v, nil
}

func TestConnection(t *testing.T) {
	apikey, err := getAPIKeyFromEnv()
	if err != nil {
		panic(err)
	}

	conn, err := NewConnectionWithAPIKey(temporalCloudAPIAddress, false, apikey)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	client := cloudservice.NewCloudServiceClient(conn)

	resp, err := client.GetUsers(context.TODO(), &cloudservice.GetUsersRequest{})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUsers())
}
