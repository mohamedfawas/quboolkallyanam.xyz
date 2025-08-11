package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

func (h *UserHandler) ConfirmAdditionalPhotoUpload(ctx context.Context,
	req *userpbv1.ConfirmAdditionalPhotoUploadRequest) (*userpbv1.ConfirmAdditionalPhotoUploadResponse, error) {

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

	additionalPhotoURL, err := h.userProfileUsecase.ConfirmAdditionalPhotoUpload(
		ctx,
		userIDUUID,
		req.ObjectKey.Value,
	)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to confirm additional photo upload", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Additional photo upload confirmed successfully")
	return &userpbv1.ConfirmAdditionalPhotoUploadResponse{
		Success:            &wrapperspb.BoolValue{Value: true},
		AdditionalPhotoUrl: &wrapperspb.StringValue{Value: additionalPhotoURL},
	}, nil
}
