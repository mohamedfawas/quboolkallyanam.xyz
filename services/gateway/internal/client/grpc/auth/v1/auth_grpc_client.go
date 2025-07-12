package v1

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
)

type authGRPCClient struct {
	conn   *grpc.ClientConn
	client authpbv1.AuthServiceClient
}

func NewAuthGRPCClient(ctx context.Context, address string, useTLS bool, tlsConfig *tls.Config) (client.AuthClient, error) {
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

	return &authGRPCClient{
		conn:   conn,
		client: authpbv1.NewAuthServiceClient(conn),
	}, nil
}

func (c *authGRPCClient) UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	grpcReq := MapUserRegisterRequest(req)
	grpcResp, err := c.client.UserRegister(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapUserRegisterResponse(grpcResp), nil
}

func (c *authGRPCClient) UserVerification(ctx context.Context, req dto.UserVerificationRequest) (*dto.UserVerificationResponse, error) {
	grpcReq := MapUserVerificationRequest(req)
	grpcResp, err := c.client.UserVerification(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapUserVerificationResponse(grpcResp), nil
}

func (c *authGRPCClient) UserLogin(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	grpcReq := MapUserLoginRequest(req)
	grpcResp, err := c.client.UserLogin(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapUserLoginResponse(grpcResp), nil
}

func (c *authGRPCClient) UserLogout(ctx context.Context, req dto.UserLogoutRequest) error {
	grpcReq := MapUserLogoutRequest(req)
	_, err := c.client.UserLogout(ctx, grpcReq)
	return err
}

func (c *authGRPCClient) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	grpcReq := MapAdminLoginRequest(req)
	grpcResp, err := c.client.AdminLogin(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapAdminLoginResponse(grpcResp), nil
}

func (c *authGRPCClient) AdminLogout(ctx context.Context, req dto.AdminLogoutRequest) error {
	grpcReq := MapAdminLogoutRequest(req)
	_, err := c.client.AdminLogout(ctx, grpcReq)
	return err
}

func (c *authGRPCClient) UserDelete(ctx context.Context, req dto.UserDeleteRequest) error {
	grpcReq := MapUserDeleteRequest(req)
	_, err := c.client.UserDelete(ctx, grpcReq)
	return err
}

func (c *authGRPCClient) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	grpcReq := MapRefreshTokenRequest(req)
	grpcResp, err := c.client.RefreshToken(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapRefreshTokenResponse(grpcResp), nil
}

func (c *authGRPCClient) Close() error {
	return c.conn.Close()
}
