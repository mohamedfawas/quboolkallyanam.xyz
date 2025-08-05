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

func (h *UserHandler) DeleteAdditionalPhoto(ctx context.Context, 
	req *userpbv1.DeleteAdditionalPhotoRequest) (*userpbv1.DeleteAdditionalPhotoResponse, error) {

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

	err = h.userProfileUsecase.DeleteAdditionalPhoto(ctx, userIDUUID, req.DisplayOrder.Value)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to delete additional photo", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Additional photo deleted successfully")
	return &userpbv1.DeleteAdditionalPhotoResponse{Success: &wrapperspb.BoolValue{Value: true}}, nil
}