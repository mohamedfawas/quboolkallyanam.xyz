package v1

import (
	"context"

	"github.com/google/uuid"
	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
	"go.uber.org/zap"
)

func (h *UserHandler) GetUserDetailsByProfileID(
	ctx context.Context,
	req *userpbv1.GetUserDetailsByProfileIDRequest,
) (*userpbv1.GetUserDetailsByProfileIDResponse, error) {

	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	requesterUUID, err := uuid.Parse(contextData.UserID)
	if err != nil {
		log.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	profile, prefs, photos, err := h.userProfileUsecase.GetUserDetailsByProfileID(ctx, requesterUUID, req.TargetProfileId, req.RequestedByAdmin)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get user details by profile id", zap.Error(err))
		}
		return nil, err
	}

	var protoPrefs *userpbv1.PartnerPreference
	if prefs != nil {
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

		protoPrefs = &userpbv1.PartnerPreference{
			MinAgeYears:               int32(prefs.MinAgeYears),
			MaxAgeYears:               int32(prefs.MaxAgeYears),
			MinHeightCm:               int32(prefs.MinHeightCm),
			MaxHeightCm:               int32(prefs.MaxHeightCm),
			AcceptPhysicallyChallenged: prefs.AcceptPhysicallyChallenged,
			PreferredCommunities:       preferredCommunities,
			PreferredMaritalStatus:     preferredMaritalStatus,
			PreferredProfessions:       preferredProfessions,
			PreferredProfessionTypes:   preferredProfessionTypes,
			PreferredEducationLevels:   preferredEducationLevels,
			PreferredHomeDistricts:     preferredHomeDistricts,
		}
	} else {
		protoPrefs = &userpbv1.PartnerPreference{}
	}

	resp := &userpbv1.GetUserDetailsByProfileIDResponse{
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
		PartnerPreferences: protoPrefs,
		AdditionalPhotoUrls: photos, // may be empty
	}

	log.Info("Successfully fetched user details by profile id")
	return resp, nil
}