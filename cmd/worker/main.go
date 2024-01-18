package main

import (
	"fmt"
	"os"

	"github.com/temporalio/cloud-samples-go/client/api"
	"github.com/temporalio/cloud-samples-go/client/temporal"
	"github.com/temporalio/cloud-samples-go/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/server/common/log"

	"go.uber.org/zap"
)

const (
	temporalCloudAPIAddress    = "saas-api.tmprl.cloud:443"
	temporalCloudAPIKeyEnvName = "TEMPORAL_CLOUD_API_KEY"

	temporalCloudNamespaceEnvName = "TEMPORAL_CLOUD_NAMESPACE"
	temporalCloudTLSCertPathEnv   = "TEMPORAL_CLOUD_TLS_CERT"
	temporalCloudTLSKeyPathEnv    = "TEMPORAL_CLOUD_TLS_KEY"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	apikey, err := getAPIKeyFromEnv()
	if err != nil {
		panic(err)
	}
	c, err := newTemporalClient(logger)
	if err != nil {
		panic(fmt.Errorf("failed to create temporal client: %+v", err))
	}
	defer c.Close()
	w := newWorker(c)

	conn, err := api.NewConnectionWithAPIKey(temporalCloudAPIAddress, false, apikey)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	workflows.Register(w, workflows.NewWorkflows(), workflows.NewActivities(conn))
	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(fmt.Errorf("failed to run worker: %+v", err))
	}
}

func newTemporalClient(logger *zap.Logger) (client.Client, error) {
	ns := os.Getenv(temporalCloudNamespaceEnvName)
	if ns == "" {
		return client.Dial(client.Options{})
	}
	return temporal.GetTemporalCloudNamespaceClient(&temporal.GetTemporalCloudNamespaceClientInput{
		Namespace:       ns,
		TLSCertFilePath: os.Getenv(temporalCloudTLSCertPathEnv),
		TLSKeyFilePath:  os.Getenv(temporalCloudTLSKeyPathEnv),
		Logger:          log.NewSdkLogger(log.NewZapLogger(logger)),
	})
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
