package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

func (h *UserHandler) GetProfilePhotoUploadURL(
	ctx context.Context,
	req *userpbv1.GetProfilePhotoUploadURLRequest) (*userpbv1.GetProfilePhotoUploadURLResponse, error) {

	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	userIDUUID, err := uuid.Parse(contextData.UserID)
	if err != nil {
		log.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	response, err := h.userProfileUsecase.GetProfilePhotoUploadURL(ctx, userIDUUID, req.ContentType.Value)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get profile photo upload URL", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Profile photo upload URL generated successfully")
	return &userpbv1.GetProfilePhotoUploadURLResponse{
		UploadUrl:        &wrapperspb.StringValue{Value: response.UploadURL},
		ObjectKey:        &wrapperspb.StringValue{Value: response.ObjectKey},
		ExpiresInSeconds: &wrapperspb.Int32Value{Value: int32(response.ExpiresInSeconds)},
	}, nil
}