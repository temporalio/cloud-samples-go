package api

import (
	"fmt"

	"go.temporal.io/cloud-sdk/cloudclient"
)

type Client struct {
	*cloudclient.Client
}

func NewConnectionWithAPIKey(apikey string) (*Client, error) {

	var cClient *cloudclient.Client
	var err error
	cClient, err = cloudclient.New(cloudclient.Options{
		APIKey: apikey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect : %v", err)
	}

	return &Client{cClient}, nil
}
