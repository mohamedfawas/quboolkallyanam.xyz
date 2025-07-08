package admin

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	adminpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/admin/v1"
)

type Client struct {
	conn   *grpc.ClientConn
	client adminpbv1.AdminServiceClient
}

func NewClient(ctx context.Context, address string, useTLS bool, tlsConfig *tls.Config) (*Client, error) {
	var creds credentials.TransportCredentials
	if useTLS {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &Client{
		conn:   conn,
		client: adminpbv1.NewAdminServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
