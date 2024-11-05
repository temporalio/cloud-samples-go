package main

import (
	"context"
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
	temporalCloudAPIKeyEnvName           = "TEMPORAL_CLOUD_API_KEY"
	temporalCloudNamespaceEnvName        = "TEMPORAL_CLOUD_NAMESPACE"
	temporalCloudNamespaceAPIKeyEnvName  = "TEMPORAL_CLOUD_NAMESPACE_API_KEY"
	temporalCloudNamespaceTLSCertPathEnv = "TEMPORAL_CLOUD_NAMESPACE_TLS_CERT"
	temporalCloudNamespaceTLSKeyPathEnv  = "TEMPORAL_CLOUD_NAMESPACE_TLS_KEY"
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

	client, err := api.NewConnectionWithAPIKey(api.TemporalCloudAPIAddress, false, apikey)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	workflows.Register(w, workflows.NewWorkflows(), workflows.NewActivities(client))
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
	var auth temporal.AuthType
	if os.Getenv(temporalCloudNamespaceTLSKeyPathEnv) != "" || os.Getenv(temporalCloudNamespaceTLSCertPathEnv) != "" {
		// if either of the TLS cert or key path is provided try to use mTLS
		auth = &temporal.MtlsAuth{
			TLSCertFilePath: os.Getenv(temporalCloudNamespaceTLSCertPathEnv),
			TLSKeyFilePath:  os.Getenv(temporalCloudNamespaceTLSKeyPathEnv),
		}
	} else if os.Getenv(temporalCloudNamespaceAPIKeyEnvName) != "" {
		// if a namespace specific API key is provided use it
		auth = &temporal.ApiKeyAuth{
			APIKey: os.Getenv(temporalCloudNamespaceAPIKeyEnvName),
		}
	} else {
		// if no specific auth is provided fallback to using the API key provided for the control plane
		auth = &temporal.ApiKeyAuth{
			APIKey: os.Getenv(temporalCloudAPIKeyEnvName),
		}
	}

	return temporal.GetTemporalCloudNamespaceClient(
		context.Background(),
		&temporal.GetTemporalCloudNamespaceClientInput{
			Namespace: ns,
			Auth:      auth,
			Logger:    log.NewSdkLogger(log.NewZapLogger(logger)),
		},
	)
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
