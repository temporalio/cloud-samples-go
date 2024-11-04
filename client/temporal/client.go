package temporal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/temporalio/cloud-samples-go/client/api"
	"github.com/temporalio/cloud-samples-go/internal/validator"
	cloudservicev1 "go.temporal.io/api/cloud/cloudservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

const (
	localTemporalHostPort = "localhost:7233"
)

type (
	ApiKeyAuth struct {
		// The api key to use for the client
		APIKey string
	}

	MtlsAuth struct {
		// The temporal cloud namespace's grpc endpoint address
		// defaults to '<namespace>.tmprl.cloud:7233'
		GRPCEndpoint string
		// The TLS cert and key file paths
		// Read more about TLS in Temporal here: https://docs.temporal.io/cloud/Certificates
		TLSCertFilePath string
		TLSKeyFilePath  string
	}

	GetTemporalCloudNamespaceClientInput struct {
		// The temporal cloud namespace to connect to (required) for e.g. "prod.a2dd6"
		Namespace string `required:"true"`

		// The auth to use for the client, defaults to local
		Auth AuthType

		// The API key to use for the client, defaults to no API key.
		APIKey string

		// The logger to use for the client, defaults to no logging
		Logger log.Logger
	}

	AuthType interface {
		getGrpcEndpoint(ctx context.Context, namespace string) (string, error)
		apply(options *client.Options) error
	}
)

func (a *ApiKeyAuth) getGrpcEndpoint(ctx context.Context, namespace string) (string, error) {
	client, err := api.NewConnectionWithAPIKey(localTemporalHostPort, false, a.APIKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cloud api connection: %w", err)
	}
	resp, err := client.CloudService().GetNamespace(context.Background(), &cloudservicev1.GetNamespaceRequest{
		Namespace: namespace,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get namespace %s: %w", namespace, err)
	}
	return resp.GetNamespace().GetEndpoints().GetGrpcAddress(), nil
}

func (a *ApiKeyAuth) apply(options *client.Options) error {
	options.Credentials = client.NewAPIKeyStaticCredentials(a.APIKey)
	return nil
}

func (a *MtlsAuth) getGrpcEndpoint(ctx context.Context, namespace string) (string, error) {
	if a.GRPCEndpoint != "" {
		return a.GRPCEndpoint, nil
	}
	return fmt.Sprintf("%s.tmprl.cloud:7233", namespace), nil
}

func (a *MtlsAuth) apply(options *client.Options) error {
	if a.TLSCertFilePath == "" || a.TLSKeyFilePath == "" {
		return fmt.Errorf("both tls cert and key file paths are required")
	}
	serverName, _, parseErr := net.SplitHostPort(a.GRPCEndpoint)
	if parseErr != nil {
		return fmt.Errorf("failed to split hostport %s: %w", a.GRPCEndpoint, parseErr)
	}
	var cert tls.Certificate
	var err error
	cert, err = tls.LoadX509KeyPair(a.TLSCertFilePath, a.TLSKeyFilePath)
	if err != nil {
		return fmt.Errorf("failed to load TLS from files: %w", err)
	}
	options.ConnectionOptions = client.ConnectionOptions{TLS: &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}}
	return nil
}

func GetTemporalCloudNamespaceClient(ctx context.Context, input *GetTemporalCloudNamespaceClientInput) (client.Client, error) {

	err := validator.ValidateStruct(input)
	if err != nil {
		return nil, err
	}
	endpoint, err := input.Auth.getGrpcEndpoint(ctx, input.Namespace)
	if err != nil {
		return nil, err
	}
	opts := client.Options{
		HostPort:  endpoint,
		Namespace: input.Namespace,
		Logger:    input.Logger,
	}
	err = input.Auth.apply(&opts)
	if err != nil {
		return nil, err
	}
	if input.Logger != nil {
		opts.Logger = input.Logger
	}
	return client.Dial(opts)
}
