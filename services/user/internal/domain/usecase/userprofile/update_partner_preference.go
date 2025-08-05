package userprofile

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) UpdateUserPartnerPreferences(
	ctx context.Context,
	userID uuid.UUID,
	operationType string,
	req entity.UpdateUserPartnerPreferencesRequest) error {

	existingProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		if err == appError.ErrUserProfileNotFound {
			return appError.ErrUserProfileNotFound
		}
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	existingPartnerPreferences, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx,
		existingProfile.ID)
	if err != nil {
		return fmt.Errorf("failed to get partner preferences: %w", err)
	}

	if operationType == constants.CreateOperationType {
		if existingPartnerPreferences != nil {
			return appError.ErrPartnerPreferencesAlreadyExists
		}
		return u.createPartnerPreferences(ctx, existingProfile.ID, req)
	}

	// operation type update, but no record exists
	if existingPartnerPreferences == nil {
		return appError.ErrPartnerPreferencesNotFound
	}

	return u.patchPartnerPreferences(ctx, existingProfile.ID, existingPartnerPreferences, req)
}

func (u *userProfileUsecase) createPartnerPreferences(
	ctx context.Context,
	userProfileID int64,
	req entity.UpdateUserPartnerPreferencesRequest) error {

	preferences := &entity.PartnerPreference{
		UserProfileID:              userProfileID,
		AcceptPhysicallyChallenged: *req.AcceptPhysicallyChallenged,
	}

	if validation.IsValidAgeRange(int(*req.MinAgeYears), int(*req.MaxAgeYears)) {
		preferences.MinAgeYears = uint16(*req.MinAgeYears)
		preferences.MaxAgeYears = uint16(*req.MaxAgeYears)
	} else {
		return appError.ErrInvalidAgeRange
	}

	if validation.IsValidHeightRange(int(*req.MinHeightCM), int(*req.MaxHeightCM)) {
		preferences.MinHeightCm = uint16(*req.MinHeightCM)
		preferences.MaxHeightCm = uint16(*req.MaxHeightCM)
	} else {
		return appError.ErrInvalidHeightRange
	}

	communities, err := validation.ParsePreferredCommunities(*req.PreferredCommunities)
	if err != nil {
		return appError.ErrInvalidCommunity
	}
	preferences.PreferredCommunities = communities

	maritalStatuses, err := validation.ParsePreferredMaritalStatuses(*req.PreferredMaritalStatus)
	if err != nil {
		return appError.ErrInvalidMaritalStatus
	}
	preferences.PreferredMaritalStatus = maritalStatuses

	professions, err := validation.ParsePreferredProfessions(*req.PreferredProfessions)
	if err != nil {
		return appError.ErrInvalidProfession
	}
	preferences.PreferredProfessions = professions

	professionTypes, err := validation.ParsePreferredProfessionTypes(*req.PreferredProfessionTypes)
	if err != nil {
		return appError.ErrInvalidProfessionType
	}
	preferences.PreferredProfessionTypes = professionTypes

	educationLevels, err := validation.ParsePreferredEducationLevels(*req.PreferredEducationLevels)
	if err != nil {
		return appError.ErrInvalidEducationLevel
	}
	preferences.PreferredEducationLevels = educationLevels

	homeDistricts, err := validation.ParsePreferredHomeDistricts(*req.PreferredHomeDistricts)
	if err != nil {
		return appError.ErrInvalidHomeDistrict
	}
	preferences.PreferredHomeDistricts = homeDistricts

	now := time.Now().UTC()
	preferences.CreatedAt = now
	preferences.UpdatedAt = now

	return u.partnerPreferencesRepository.CreatePartnerPreferences(ctx, preferences)
}

func (u *userProfileUsecase) patchPartnerPreferences(
	ctx context.Context,
	userProfileID int64,
	existingPartnerPreferences *entity.PartnerPreference,
	req entity.UpdateUserPartnerPreferencesRequest) error {

	patch := make(map[string]interface{})

	// Handle age range updates
	if req.MinAgeYears != nil && req.MaxAgeYears != nil {
		if validation.IsValidAgeRange(int(*req.MinAgeYears), int(*req.MaxAgeYears)) {
			patch["min_age_years"] = uint16(*req.MinAgeYears)
			patch["max_age_years"] = uint16(*req.MaxAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	} else if req.MinAgeYears != nil {
		if validation.IsValidAge(int(*req.MinAgeYears)) && int(*req.MinAgeYears) <= int(existingPartnerPreferences.MaxAgeYears) {
			patch["min_age_years"] = uint16(*req.MinAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	} else if req.MaxAgeYears != nil {
		if validation.IsValidAge(int(*req.MaxAgeYears)) && int(*req.MaxAgeYears) >= int(existingPartnerPreferences.MinAgeYears) {
			patch["max_age_years"] = uint16(*req.MaxAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	}

	// Handle height range updates
	if req.MinHeightCM != nil && req.MaxHeightCM != nil {
		if validation.IsValidHeightRange(int(*req.MinHeightCM), int(*req.MaxHeightCM)) {
			patch["min_height_cm"] = uint16(*req.MinHeightCM)
			patch["max_height_cm"] = uint16(*req.MaxHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	} else if req.MinHeightCM != nil {
		if validation.IsValidHeight(int(*req.MinHeightCM)) && int(*req.MinHeightCM) <= int(existingPartnerPreferences.MaxHeightCm) {
			patch["min_height_cm"] = uint16(*req.MinHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	} else if req.MaxHeightCM != nil {
		if validation.IsValidHeight(int(*req.MaxHeightCM)) && int(*req.MaxHeightCM) >= int(existingPartnerPreferences.MinHeightCm) {
			patch["max_height_cm"] = uint16(*req.MaxHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	}

	// Handle physically challenged preference
	if req.AcceptPhysicallyChallenged != nil {
		patch["accept_physically_challenged"] = *req.AcceptPhysicallyChallenged
	}

	if len(*req.PreferredCommunities) > 0 {
		communities, err := validation.ParsePreferredCommunities(*req.PreferredCommunities)
		if err != nil {
			return appError.ErrInvalidCommunity
		}
		patch["preferred_communities"] = communities
	}

	if len(*req.PreferredMaritalStatus) > 0 {
		maritalStatuses, err := validation.ParsePreferredMaritalStatuses(*req.PreferredMaritalStatus)
		if err != nil {
			return appError.ErrInvalidMaritalStatus
		}
		patch["preferred_marital_status"] = maritalStatuses
	}

	if len(*req.PreferredProfessions) > 0 {
		professions, err := validation.ParsePreferredProfessions(*req.PreferredProfessions)
		if err != nil {
			return appError.ErrInvalidProfession
		}
		patch["preferred_professions"] = professions
	}

	if len(*req.PreferredProfessionTypes) > 0 {
		professionTypes, err := validation.ParsePreferredProfessionTypes(*req.PreferredProfessionTypes)
		if err != nil {
			return appError.ErrInvalidProfessionType
		}
		patch["preferred_profession_types"] = professionTypes
	}

	if len(*req.PreferredEducationLevels) > 0 {
		educationLevels, err := validation.ParsePreferredEducationLevels(*req.PreferredEducationLevels)
		if err != nil {
			return appError.ErrInvalidEducationLevel
		}
		patch["preferred_education_levels"] = educationLevels
	}

	if len(*req.PreferredHomeDistricts) > 0 {
		homeDistricts, err := validation.ParsePreferredHomeDistricts(*req.PreferredHomeDistricts)
		if err != nil {
			return appError.ErrInvalidHomeDistrict
		}
		patch["preferred_home_districts"] = homeDistricts
	}

	if len(patch) > 0 {
		patch["updated_at"] = time.Now().UTC()
	}

	return u.partnerPreferencesRepository.PatchPartnerPreferences(ctx, userProfileID, patch)
}
