package main

import (
	"fmt"
	"os"

	"github.com/temporalio/cloud-samples-go/client/temporalcloud"
	"github.com/temporalio/cloud-samples-go/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	temporalCloudAPIAddress    = "saas-api.tmprl.cloud:443"
	temporalCloudAPIKeyEnvName = "TEMPORAL_CLOUD_API_KEY"
)

func main() {
	apikey, err := getAPIKeyFromEnv()
	if err != nil {
		panic(err)
	}
	c, err := newLocalTemporalClient()
	if err != nil {
		panic(fmt.Errorf("failed to create temporal client: %+v", err))
	}
	defer c.Close()
	w := newWorker(c)

	conn, err := temporalcloud.NewConnectionWithAPIKey(temporalCloudAPIAddress, false, apikey)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	workflows.Register(w, workflows.NewWorkflows(), workflows.NewActivities(conn))
	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(fmt.Errorf("failed to run worker: %+v", err))
	}
}

func newLocalTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{})
}

func newWorker(client client.Client) worker.Worker {
	wo := worker.Options{
		MaxConcurrentActivityTaskPollers: 10,
		MaxConcurrentWorkflowTaskPollers: 10,
	}
	return worker.New(client, "demo", wo)
}

func getAPIKeyFromEnv() (string, error) {
	v := os.Getenv(temporalCloudAPIKeyEnvName)
	if v == "" {
		return "", fmt.Errorf("apikey not provided, set environment variable '%s' with apikey you want to use", temporalCloudAPIKeyEnvName)
	}
	return v, nil
}
