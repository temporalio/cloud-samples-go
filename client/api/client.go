package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func NewConnectionWithAPIKey(addrStr string, allowInsecure bool, apiKey string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return NewConnection(
		addrStr,
		allowInsecure,
		append(opts, grpc.WithPerRPCCredentials(NewAPIKeyRPCCredential(apiKey, allowInsecure)))...,
	)
}

func NewConnection(addrStr string, allowInsecure bool, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	addr, err := url.Parse(addrStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse server address: %s", err)
	}
	defaultOpts, err := defaultDialOptions(addr, allowInsecure)
	if err != nil {
		return nil, fmt.Errorf("failed to generate default dial options: %s", err)
	}

	conn, err := grpc.Dial(
		addr.String(),
		append(defaultOpts, opts...)...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial `%s`: %v", addr.String(), err)
	}
	return conn, nil
}

func defaultDialOptions(addr *url.URL, allowInsecure bool) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	transport := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: addr.Hostname(),
	})
	if allowInsecure {
		transport = insecure.NewCredentials()
	}

	opts = append(opts, grpc.WithTransportCredentials(transport))
	opts = append(opts, grpc.WithUnaryInterceptor(setAPIVersionInterceptor))
	return opts, nil
}

func setAPIVersionInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx = metadata.AppendToOutgoingContext(ctx, TemporalCloudAPIVersionHeader, strings.TrimSpace(TemporalCloudAPIVersion))
	return invoker(ctx, method, req, reply, cc, opts...)
}
