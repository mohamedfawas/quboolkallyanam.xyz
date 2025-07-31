package user

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
)

type userGRPCClient struct {
	conn   *grpc.ClientConn
	client userpbv1.UserServiceClient
}

func NewUserGRPCClient(
	ctx context.Context,
	address string,
	useTLS bool,
	tlsConfig *tls.Config) (client.UserClient, error) {
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

	cc, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client to %s: %w", address, err)
	}

	return &userGRPCClient{
		conn:   cc,
		client: userpbv1.NewUserServiceClient(cc),
	}, nil
}

func (c *userGRPCClient) UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapUpdateUserProfileRequest(req)
	_, err := c.client.UpdateUserProfile(ctx, grpcReq)
	if err != nil {
		return err
	}
	return nil
}

func (c *userGRPCClient) UpdateUserPartnerPreferences(ctx context.Context, operationType string, req dto.PartnerPreferencePatchRequest) error {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapUpdateUserPartnerPreferencesRequest(operationType, req)
	_, err := c.client.UpdateUserPartnerPreferences(ctx, grpcReq)
	if err != nil {
		return err
	}
	return nil
}

func (c *userGRPCClient) RecordMatchAction(ctx context.Context, req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapRecordMatchActionRequest(req)
	resp, err := c.client.RecordMatchAction(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapRecordMatchActionResponse(resp), nil
}

func (c *userGRPCClient) GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapGetMatchRecommendationsRequest(req)
	resp, err := c.client.GetMatchRecommendations(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetMatchRecommendationsResponse(resp), nil
}

func (c *userGRPCClient) GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	ctx = contextutils.SetUserContext(ctx, userID)
	grpcReq := MapGetProfilesByMatchActionRequest(req)
	resp, err := c.client.GetProfilesByMatchAction(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetProfilesByMatchActionResponse(resp), nil
}

func (c *userGRPCClient) Close() error {
	return c.conn.Close()
}
