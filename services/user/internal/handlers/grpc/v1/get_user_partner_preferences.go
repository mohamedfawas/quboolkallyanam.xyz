package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

func (h *UserHandler) GetUserPartnerPreferences(
	ctx context.Context,
	req *userpbv1.GetUserPartnerPreferencesRequest,
) (*userpbv1.GetUserPartnerPreferencesResponse, error) {

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

	prefs, err := h.userProfileUsecase.GetUserPartnerPreferences(ctx, userIDUUID)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get partner preferences", zap.Error(err))
		}
		return nil, err
	}

	preferredCommunities := make([]string, len(prefs.PreferredCommunities))
	for i, v := range prefs.PreferredCommunities {
		preferredCommunities[i] = string(v)
	}
	preferredMaritalStatus := make([]string, len(prefs.PreferredMaritalStatus))
	for i, v := range prefs.PreferredMaritalStatus {
		preferredMaritalStatus[i] = string(v)
	}
	preferredProfessions := make([]string, len(prefs.PreferredProfessions))
	for i, v := range prefs.PreferredProfessions {
		preferredProfessions[i] = string(v)
	}
	preferredProfessionTypes := make([]string, len(prefs.PreferredProfessionTypes))
	for i, v := range prefs.PreferredProfessionTypes {
		preferredProfessionTypes[i] = string(v)
	}
	preferredEducationLevels := make([]string, len(prefs.PreferredEducationLevels))
	for i, v := range prefs.PreferredEducationLevels {
		preferredEducationLevels[i] = string(v)
	}
	preferredHomeDistricts := make([]string, len(prefs.PreferredHomeDistricts))
	for i, v := range prefs.PreferredHomeDistricts {
		preferredHomeDistricts[i] = string(v)
	}

	resp := &userpbv1.GetUserPartnerPreferencesResponse{
		PartnerPreferences: &userpbv1.PartnerPreference{
			MinAgeYears:                int32(prefs.MinAgeYears),
			MaxAgeYears:                int32(prefs.MaxAgeYears),
			MinHeightCm:                int32(prefs.MinHeightCm),
			MaxHeightCm:                int32(prefs.MaxHeightCm),
			AcceptPhysicallyChallenged: prefs.AcceptPhysicallyChallenged,
			PreferredCommunities:       preferredCommunities,
			PreferredMaritalStatus:     preferredMaritalStatus,
			PreferredProfessions:       preferredProfessions,
			PreferredProfessionTypes:   preferredProfessionTypes,
			PreferredEducationLevels:   preferredEducationLevels,
			PreferredHomeDistricts:     preferredHomeDistricts,
		},
	}

	log.Info("Successfully fetched partner preferences")
	return resp, nil
}