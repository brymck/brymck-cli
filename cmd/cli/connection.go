package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", t.token),
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func getConnection(addr string) (*grpc.ClientConn, error) {
	pool, _ := x509.SystemCertPool()
	ce := credentials.NewClientTLSFromCert(pool, "")

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(ce),
		grpc.WithPerRPCCredentials(tokenAuth{token: os.Getenv("BRYMCK_ID_TOKEN")}),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
