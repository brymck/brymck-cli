package connections

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const ServiceAddressTemplate = "%s-4tt23pryoq-an.a.run.app:443"

type Connection struct {
	Context              context.Context
	GrpcClientConnection *grpc.ClientConn
	cancel               context.CancelFunc
}

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", t.token),
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func getGrpcClientConnection(addr string) (*grpc.ClientConn, error) {
	if strings.Contains(addr, "localhost") {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return conn, nil
	} else {

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
}

func (conn *Connection) Close() {
	conn.GrpcClientConnection.Close()
	conn.cancel()
}

func NewConnection(serviceName string, addr string) (*Connection, error) {
	if addr == "" {
		addr = fmt.Sprintf(ServiceAddressTemplate, serviceName)
	}
	conn, err := getGrpcClientConnection(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	return &Connection{Context: ctx, GrpcClientConnection: conn, cancel: cancel}, nil
}
