package temporal

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/temporalio/cloud-samples-go/internal/validator"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

const (
	localTemporalHostPort = "localhost:7233"
)

type (
	GetTemporalCloudNamespaceClientInput struct {
		// The temporal cloud namespace to connect to (required) for e.g. "prod.a2dd6"
		Namespace string `required:"true"`
		// The temporal cloud namespace's grpc endpoint address, defaults to '<namespace>.tmprl.cloud:7233'
		GRPCEndpoint string

		// The TLS cert and key file paths
		// Read more about TLS in Temporal here: https://docs.temporal.io/cloud/certificates
		TLSCertFilePath string `required:"true"`
		TLSKeyFilePath  string `required:"true"`

		// The logger to use for the client, defaults to no logging
		Logger log.Logger
	}
)

func GetTemporalCloudNamespaceClient(input *GetTemporalCloudNamespaceClientInput) (client.Client, error) {

	err := validator.ValidateStruct(input)
	if err != nil {
		return nil, err
	}
	if input.GRPCEndpoint == "" {
		input.GRPCEndpoint = fmt.Sprintf("%s.tmprl.cloud:7233", input.Namespace)
	}
	tlsConfig, err := getTLSConfig(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	opts := client.Options{
		HostPort:          input.GRPCEndpoint,
		Namespace:         input.Namespace,
		ConnectionOptions: client.ConnectionOptions{TLS: tlsConfig},
		Logger:            input.Logger,
	}
	if input.Logger != nil {
		opts.Logger = input.Logger
	}
	return client.Dial(opts)
}

func getTLSConfig(input *GetTemporalCloudNamespaceClientInput) (*tls.Config, error) {
	if input.TLSCertFilePath == "" || input.TLSKeyFilePath == "" {
		return nil, nil
	}
	serverName, _, parseErr := net.SplitHostPort(input.GRPCEndpoint)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to split hostport %s: %w", input.GRPCEndpoint, parseErr)
	}
	var cert tls.Certificate
	var err error
	cert, err = tls.LoadX509KeyPair(input.TLSCertFilePath, input.TLSKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS from files: %w", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}, nil
}
