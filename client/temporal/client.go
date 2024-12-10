package temporal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/temporalio/cloud-samples-go/internal/validator"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	ApiKeyAuth struct {
		// The func that returns the api key to use
		GetApiKeyCallback func(context.Context) (string, error) `required:"true"`
		// The temporal cloud namespace's grpc endpoint address to connect to
		GrpcAddress string `required:"false"`
	}

	MtlsAuth struct {
		// The temporal cloud namespace's grpc endpoint address
		// defaults to '<namespace>.tmprl.cloud:7233'
		GRPCEndpoint string
		// The TLS cert and key file paths
		// Read more about TLS in Temporal here: https://docs.temporal.io/cloud/Certificates
		TLSCertFilePath string `required:"true"`
		TLSKeyFilePath  string `required:"true"`
	}

	GetTemporalCloudNamespaceClientInput struct {
		// The temporal cloud namespace to connect to (required) for e.g. "prod.a2dd6"
		Namespace string `required:"true"`

		// The auth to use for the client, defaults to local
		Auth AuthType

		// The logger to use for the client, defaults to no logging
		Logger log.Logger
	}

	AuthType interface {
		apply(options *client.Options) error
	}
)

func (a *ApiKeyAuth) apply(options *client.Options) error {
	options.HostPort = a.GrpcAddress
	options.Credentials = client.NewAPIKeyDynamicCredentials(a.GetApiKeyCallback)
	options.ConnectionOptions = client.ConnectionOptions{
		TLS: &tls.Config{},
		DialOptions: []grpc.DialOption{
			grpc.WithUnaryInterceptor(
				func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
					return invoker(
						metadata.AppendToOutgoingContext(ctx, "temporal-namespace", options.Namespace),
						method,
						req,
						reply,
						cc,
						opts...,
					)
				},
			),
		},
	}
	return nil
}

func (a *MtlsAuth) apply(options *client.Options) error {
	endpoint := a.GRPCEndpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s.tmprl.cloud:7233", options.Namespace)
	}
	if a.TLSCertFilePath == "" || a.TLSKeyFilePath == "" {
		return fmt.Errorf("both tls cert and key file paths are required")
	}
	serverName, _, parseErr := net.SplitHostPort(endpoint)
	if parseErr != nil {
		return fmt.Errorf("failed to split hostport %s: %w", endpoint, parseErr)
	}
	options.HostPort = endpoint
	options.ConnectionOptions = client.ConnectionOptions{TLS: &tls.Config{
		ServerName: serverName,
	}}
	a.setClientCertAutoRefresh(options.ConnectionOptions.TLS)
	return nil
}

func (a *MtlsAuth) setClientCertAutoRefresh(tlsConfig *tls.Config) error {
	lastModifiedTime := time.Time{}
	var clientCert *tls.Certificate
	var mu sync.Mutex

	tlsConfig.GetClientCertificate = func(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		mu.Lock()
		defer mu.Unlock()

		fileInfo, err := os.Stat(a.TLSCertFilePath)
		if err != nil {
			return clientCert, fmt.Errorf("stat error for client tls cert: %w", err)
		}

		newModifiedTime := fileInfo.ModTime()
		if newModifiedTime.Equal(lastModifiedTime) {
			return clientCert, nil
		}

		newClientCert, err := tls.LoadX509KeyPair(a.TLSCertFilePath, a.TLSKeyFilePath)
		if err != nil {
			return clientCert, fmt.Errorf("failed to load TLS from files: %w", err)
		}

		lastModifiedTime = newModifiedTime
		clientCert = &newClientCert
		return clientCert, nil
	}

	return nil
}

func GetTemporalCloudNamespaceClient(ctx context.Context, input *GetTemporalCloudNamespaceClientInput) (client.Client, error) {
	err := validator.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	opts := client.Options{
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
