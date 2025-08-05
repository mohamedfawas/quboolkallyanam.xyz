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

func (h *UserHandler) GetAdditionalPhotoUploadURL(ctx context.Context,
	req *userpbv1.GetAdditionalPhotoUploadURLRequest) (*userpbv1.GetAdditionalPhotoUploadURLResponse, error) {

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

	response, err := h.userProfileUsecase.GetAdditionalPhotoUploadURL(ctx, userIDUUID, int32(req.DisplayOrder.Value), req.ContentType.Value)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get additional photo upload URL", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Additional photo upload URL generated successfully")

	return &userpbv1.GetAdditionalPhotoUploadURLResponse{
		UploadUrl:        &wrapperspb.StringValue{Value: response.UploadURL},
		ObjectKey:        &wrapperspb.StringValue{Value: response.ObjectKey},
		ExpiresInSeconds: &wrapperspb.UInt32Value{Value: response.ExpiresInSeconds},
	}, nil
}
