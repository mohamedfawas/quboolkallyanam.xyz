package user

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

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

///////// USER PROFILE MANAGEMENT //////////
func (c *userGRPCClient) UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return err
	}

	grpcReq := MapUpdateUserProfileRequest(req)
	_, err = c.client.UpdateUserProfile(ctx, grpcReq)
	if err != nil {
		return err
	}
	return nil
}

func (c *userGRPCClient) GetUserProfile(ctx context.Context) (*dto.UserProfileRecommendation, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.GetUserProfile(ctx, &userpbv1.GetUserProfileRequest{})
	if err != nil {
		return nil, err
	}
	return MapGetUserProfileResponse(resp), nil
}

func (c *userGRPCClient) GetProfilePhotoUploadURL(ctx context.Context, req dto.GetProfilePhotoUploadURLRequest) (*dto.GetProfilePhotoUploadURLResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapGetProfilePhotoUploadURLRequest(req)
	resp, err := c.client.GetProfilePhotoUploadURL(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetProfilePhotoUploadURLResponse(resp), nil
}


func (c *userGRPCClient) ConfirmProfilePhotoUpload(ctx context.Context, req dto.ConfirmProfilePhotoUploadRequest) (*dto.ConfirmProfilePhotoUploadResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapConfirmProfilePhotoUploadRequest(req)
	resp, err := c.client.ConfirmProfilePhotoUpload(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapConfirmProfilePhotoUploadResponse(resp), nil
}

func (c *userGRPCClient) DeleteProfilePhoto(ctx context.Context, req dto.DeleteProfilePhotoRequest) (*dto.DeleteProfilePhotoResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapDeleteProfilePhotoRequest(req)
	resp, err := c.client.DeleteProfilePhoto(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapDeleteProfilePhotoResponse(resp), nil
}

///////////////////// User Additional Photo Management //////////////////

func (c *userGRPCClient) GetAdditionalPhotoUploadURL(ctx context.Context, req dto.GetAdditionalPhotoUploadURLRequest) (*dto.GetAdditionalPhotoUploadURLResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapGetAdditionalPhotoUploadURLRequest(req)
	resp, err := c.client.GetAdditionalPhotoUploadURL(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetAdditionalPhotoUploadURLResponse(resp), nil
}


func (c *userGRPCClient) ConfirmAdditionalPhotoUpload(ctx context.Context, req dto.ConfirmAdditionalPhotoUploadRequest) (*dto.ConfirmAdditionalPhotoUploadResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapConfirmAdditionalPhotoUploadRequest(req)
	resp, err := c.client.ConfirmAdditionalPhotoUpload(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapConfirmAdditionalPhotoUploadResponse(resp), nil
}

func (c *userGRPCClient) DeleteAdditionalPhoto(ctx context.Context, req dto.DeleteAdditionalPhotoRequest) (*dto.DeleteAdditionalPhotoResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}
	
	grpcReq := MapDeleteAdditionalPhotoRequest(req)
	resp, err := c.client.DeleteAdditionalPhoto(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapDeleteAdditionalPhotoResponse(resp), nil
}

func (c *userGRPCClient) GetAdditionalPhotos(ctx context.Context) (*dto.GetAdditionalPhotosResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.GetAdditionalPhotos(ctx, &userpbv1.GetAdditionalPhotosRequest{})
	if err != nil {
		return nil, err
	}
	return MapGetAdditionalPhotosResponse(resp), nil
}

//////////////// PARTNER PREFERENCES MANAGEMENT //////////////////
func (c *userGRPCClient) UpdateUserPartnerPreferences(
	ctx context.Context, 
	operationType string, 
	req dto.UpdatePartnerPreferenceRequest) error {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return err
	}

	grpcReq := MapUpdateUserPartnerPreferencesRequest(operationType, req)
	_, err = c.client.UpdateUserPartnerPreferences(ctx, grpcReq)
	if err != nil {
		return err
	}
	return nil
}

/////////////////// MATCH MAKING ///////////////////
func (c *userGRPCClient) RecordMatchAction(ctx context.Context, req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapRecordMatchActionRequest(req)
	resp, err := c.client.RecordMatchAction(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapRecordMatchActionResponse(resp), nil
}

func (c *userGRPCClient) GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapGetMatchRecommendationsRequest(req)
	resp, err := c.client.GetMatchRecommendations(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetMatchRecommendationsResponse(resp), nil
}

func (c *userGRPCClient) GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

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
