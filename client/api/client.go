package api

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

type Client struct {
	client.CloudOperationsClient
}

var (
	_ client.CloudOperationsClient = &Client{}

	TemporalCloudAPIVersion = "2024-05-13-00"
)

func NewConnectionWithAPIKey(addrStr string, allowInsecure bool, apiKey string) (*Client, error) {

	var cClient client.CloudOperationsClient
	var err error
	cClient, err = client.DialCloudOperationsClient(context.Background(), client.CloudOperationsClientOptions{
		Version:     TemporalCloudAPIVersion,
		Credentials: client.NewAPIKeyStaticCredentials(apiKey),
		DisableTLS:  allowInsecure,
		HostPort:    addrStr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect `%s`: %v", client.DefaultHostPort, err)
	}

	return &Client{cClient}, nil
}
