package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
	"go.uber.org/zap"
)

func (h *UserHandler) GetUserProfile(
	ctx context.Context,
	req *userpbv1.GetUserProfileRequest,
) (*userpbv1.GetUserProfileResponse, error) {

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

	profile, err := h.userProfileUsecase.GetUserProfile(ctx, userIDUUID)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get user profile", zap.Error(err))
		}
		return nil, err
	}

	resp := &userpbv1.GetUserProfileResponse{
		Profile: &userpbv1.UserProfileRecommendation{
			Id:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureUrl: pointerutil.GetStringValue(profile.ProfilePictureURL),
			Age:               profile.Age,
			HeightCm:          profile.HeightCm,
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		},
	}

	log.Info("Successfully fetched user profile")
	return resp, nil
}