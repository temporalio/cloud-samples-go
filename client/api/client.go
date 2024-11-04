package api

import (
	"context"
	"fmt"
	"time"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
)

type Client struct {
	client.CloudOperationsClient
}

var (
	_ client.CloudOperationsClient = &Client{}

	TemporalCloudAPIAddress = "saas-api.tmprl.cloud:443"
	TemporalCloudAPIVersion = "2024-10-01-00"
)

func NewConnectionWithAPIKey(addrStr string, allowInsecure bool, apiKey string) (*Client, error) {

	var cClient client.CloudOperationsClient
	var err error
	cClient, err = client.DialCloudOperationsClient(context.Background(), client.CloudOperationsClientOptions{
		Version:     TemporalCloudAPIVersion,
		Credentials: client.NewAPIKeyStaticCredentials(apiKey),
		DisableTLS:  allowInsecure,
		HostPort:    addrStr,
		ConnectionOptions: client.ConnectionOptions{
			DialOptions: []grpc.DialOption{
				grpc.WithChainUnaryInterceptor(
					grpcretry.UnaryClientInterceptor(
						grpcretry.WithBackoff(
							grpcretry.BackoffExponentialWithJitter(250*time.Millisecond, 0.1),
						),
						grpcretry.WithMax(5),
					),
				),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect `%s`: %v", client.DefaultHostPort, err)
	}

	return &Client{cClient}, nil
}
