package v1

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
)

type authGRPCClient struct {
	conn   *grpc.ClientConn
	client authpbv1.AuthServiceClient
}

func NewAuthGRPCClient(
	ctx context.Context,
	address string,
	useTLS bool,
	tlsConfig *tls.Config) (client.AuthClient, error) {
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var opts []grpc.DialOption
	if useTLS {
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := insecure.NewCredentials()
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	cc, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client to %s: %w", address, err)
	}

	return &authGRPCClient{
		conn:   cc,
		client: authpbv1.NewAuthServiceClient(cc),
	}, nil
}

func (c *authGRPCClient) UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {

	grpcReq := MapUserRegisterRequest(req)
	grpcResp, err := c.client.UserRegister(ctx, grpcReq)
	if err != nil {
		log.Printf("UserRegister error in auth grpc client: %v", err)
		return nil, err
	}
	return MapUserRegisterResponse(grpcResp), nil
}

func (c *authGRPCClient) UserVerification(ctx context.Context, req dto.UserVerificationRequest) (*dto.UserVerificationResponse, error) {

	grpcReq := MapUserVerificationRequest(req)
	grpcResp, err := c.client.UserVerification(ctx, grpcReq)
	if err != nil {
		log.Printf("UserVerification error in auth grpc client: %v", err)
		return nil, err
	}

	return MapUserVerificationResponse(grpcResp), nil
}

func (c *authGRPCClient) UserLogin(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	grpcReq := MapUserLoginRequest(req)
	grpcResp, err := c.client.UserLogin(ctx, grpcReq)
	if err != nil {
		log.Printf("UserLogin error in auth grpc client: %v", err)
		return nil, err
	}
	return MapUserLoginResponse(grpcResp), nil
}

func (c *authGRPCClient) UserLogout(ctx context.Context, accessToken string) error {
	grpcReq := MapUserLogoutRequest(accessToken)
	_, err := c.client.UserLogout(ctx, grpcReq)
	if err != nil {
		log.Printf("UserLogout error in auth grpc client: %v", err)
		return err
	}
	return nil
}

func (c *authGRPCClient) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	grpcReq := MapAdminLoginRequest(req)
	grpcResp, err := c.client.AdminLogin(ctx, grpcReq)
	if err != nil {
		log.Printf("AdminLogin error in auth grpc client: %v", err)
		return nil, err
	}
	return MapAdminLoginResponse(grpcResp), nil
}

func (c *authGRPCClient) AdminLogout(ctx context.Context, accessToken string) error {
	grpcReq := MapAdminLogoutRequest(accessToken)
	_, err := c.client.AdminLogout(ctx, grpcReq)
	if err != nil {
		log.Printf("AdminLogout error in auth grpc client: %v", err)
		return err
	}
	return nil
}

func (c *authGRPCClient) UserDelete(ctx context.Context, req dto.UserDeleteRequest) error {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapUserDeleteRequest(req)
	_, err := c.client.UserDelete(ctx, grpcReq)
	if err != nil {
		log.Printf("UserDelete error in auth grpc client: %v", err)
		return err
	}
	return nil
}

func (c *authGRPCClient) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	grpcReq := MapRefreshTokenRequest(req)
	grpcResp, err := c.client.RefreshToken(ctx, grpcReq)
	if err != nil {
		log.Printf("RefreshToken error in auth grpc client: %v", err)
		return nil, err
	}
	return MapRefreshTokenResponse(grpcResp), nil
}

func (c *authGRPCClient) Close() error {
	return c.conn.Close()
}
