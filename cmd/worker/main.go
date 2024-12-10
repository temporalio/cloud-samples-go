package main

import (
	"context"
	"fmt"
	"os"

	"github.com/temporalio/cloud-samples-go/client/api"
	"github.com/temporalio/cloud-samples-go/client/temporal"
	"github.com/temporalio/cloud-samples-go/workflows"
	cloudservicev1 "go.temporal.io/api/cloud/cloudservice/v1"
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
	ctx := context.Background()
	c, err := newTemporalClient(ctx, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create temporal client: %+v", err))
	}
	defer c.Close()
	w := newWorker(c)

	client, err := api.NewConnectionWithAPIKey(api.TemporalCloudAPIAddress, false, getAPIKeyFromEnv)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	workflows.Register(w, workflows.NewWorkflows(), workflows.NewActivities(client))
	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(fmt.Errorf("failed to run worker: %+v", err))
	}
}

func newTemporalClient(ctx context.Context, logger *zap.Logger) (client.Client, error) {
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
	} else {
		// fetch the grpc address using the namespace API key
		grpcAddress, err := fetchGrpcAddressUsingAPIKey(ctx, ns, getNamespaceAPIKeyFromEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch grpc address using namespace api key: %w", err)
		}
		auth = &temporal.ApiKeyAuth{
			GetApiKeyCallback: getNamespaceAPIKeyFromEnv,
			GrpcAddress:       grpcAddress,
		}
	}

	return temporal.GetTemporalCloudNamespaceClient(
		ctx,
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

func getAPIKeyFromEnv(ctx context.Context) (string, error) {
	v := os.Getenv(temporalCloudAPIKeyEnvName)
	if v == "" {
		// if no API key is provided return an error
		return "", fmt.Errorf("apikey not provided, set environment variable '%s' with apikey you want to use", temporalCloudAPIKeyEnvName)
	}
	return v, nil
}

func getNamespaceAPIKeyFromEnv(ctx context.Context) (string, error) {
	v := os.Getenv(temporalCloudNamespaceAPIKeyEnvName)
	if v == "" {
		// fallback to using the control plane API key if no namespace specific API key is provided
		v, _ = getAPIKeyFromEnv(ctx)
	}
	if v == "" {
		// if no API key is provided return an error
		return "", fmt.Errorf("namespace apikey not provided, set environment variable '%s' with apikey you want to use", temporalCloudNamespaceAPIKeyEnvName)
	}
	return v, nil
}

func fetchGrpcAddressUsingAPIKey(
	ctx context.Context,
	namespace string,
	getAPIKeyFunc func(context.Context) (string, error),
) (string, error) {
	c, err := api.NewConnectionWithAPIKey(api.TemporalCloudAPIAddress, false, getAPIKeyFunc)
	if err != nil {
		return "", fmt.Errorf("failed to create cloud api connection: %w", err)
	}
	resp, err := c.CloudService().GetNamespace(ctx, &cloudservicev1.GetNamespaceRequest{
		Namespace: namespace,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get namespace %q: %w", namespace, err)
	}
	if resp.GetNamespace().GetEndpoints().GetGrpcAddress() == "" {
		return "", fmt.Errorf("namespace %q has no grpc address", namespace)
	}
	return resp.GetNamespace().GetEndpoints().GetGrpcAddress(), nil
}
