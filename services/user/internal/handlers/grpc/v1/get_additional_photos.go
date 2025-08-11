package v1

import (
    "context"

    "github.com/google/uuid"
    userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
    "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
    "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
    "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
    "go.uber.org/zap"
)

func (h *UserHandler) GetAdditionalPhotos(
    ctx context.Context,
    req *userpbv1.GetAdditionalPhotosRequest,
) (*userpbv1.GetAdditionalPhotosResponse, error) {

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

    urls, err := h.userProfileUsecase.GetAdditionalPhotos(ctx, userIDUUID)
    if err != nil {
        if !apperrors.IsAppError(err) {
            log.Error("Failed to get additional photos", zap.Error(err))
        }
        return nil, err
    }

    log.Info("Fetched additional photos successfully")
    return &userpbv1.GetAdditionalPhotosResponse{
        AdditionalPhotoUrls: urls,
    }, nil
}