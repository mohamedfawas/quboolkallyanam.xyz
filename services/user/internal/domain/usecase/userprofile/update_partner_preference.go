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
	"github.com/lib/pq"
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
		UserProfileID: userProfileID,
	}

	if req.AcceptPhysicallyChallenged != nil {
		preferences.AcceptPhysicallyChallenged = *req.AcceptPhysicallyChallenged
	} else {
		preferences.AcceptPhysicallyChallenged = false
	}

	if req.MinAgeYears != nil && req.MaxAgeYears != nil {
		if validation.IsValidAgeRange(int(*req.MinAgeYears), int(*req.MaxAgeYears)) {
			preferences.MinAgeYears = int16(*req.MinAgeYears)
			preferences.MaxAgeYears = int16(*req.MaxAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	} else {
		preferences.MinAgeYears = 18
		preferences.MaxAgeYears = 100
	}

	if req.MinHeightCM != nil && req.MaxHeightCM != nil {
		if validation.IsValidHeightRange(int(*req.MinHeightCM), int(*req.MaxHeightCM)) {
			preferences.MinHeightCm = int16(*req.MinHeightCM)
			preferences.MaxHeightCm = int16(*req.MaxHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	} else {
		preferences.MinHeightCm = 130
		preferences.MaxHeightCm = 220
	}

	// // Preferred communities
	communitiesIn := []string{"any"}
	if req.PreferredCommunities != nil && len(*req.PreferredCommunities) > 0 {
		communitiesIn = *req.PreferredCommunities
	}
	communitiesTyped, err := validation.ParsePreferredCommunities(communitiesIn)
	if err != nil {
		return appError.ErrInvalidCommunity
	}
	preferences.PreferredCommunities = pq.StringArray(validation.CommunitiesToStrings(communitiesTyped))

	// Marital status
	maritalIn := []string{"any"}
	if req.PreferredMaritalStatus != nil && len(*req.PreferredMaritalStatus) > 0 {
		maritalIn = *req.PreferredMaritalStatus
	}
	maritalTyped, err := validation.ParsePreferredMaritalStatuses(maritalIn)
	if err != nil {
		return appError.ErrInvalidMaritalStatus
	}
	preferences.PreferredMaritalStatus = pq.StringArray(validation.MaritalStatusesToStrings(maritalTyped))


	// Professions
	professionsIn := []string{"any"}
	if req.PreferredProfessions != nil && len(*req.PreferredProfessions) > 0 {
		professionsIn = *req.PreferredProfessions
	}
	professionsTyped, err := validation.ParsePreferredProfessions(professionsIn)
	if err != nil {
		return appError.ErrInvalidProfession
	}
	preferences.PreferredProfessions = pq.StringArray(validation.ProfessionsToStrings(professionsTyped))


	// Profession types
	profTypesIn := []string{"any"}
	if req.PreferredProfessionTypes != nil && len(*req.PreferredProfessionTypes) > 0 {
		profTypesIn = *req.PreferredProfessionTypes
	}
	professionTypesTyped, err := validation.ParsePreferredProfessionTypes(profTypesIn)
	if err != nil {
		return appError.ErrInvalidProfessionType
	}
	preferences.PreferredProfessionTypes = pq.StringArray(validation.ProfessionTypesToStrings(professionTypesTyped))


	// Education levels
	eduIn := []string{"any"}
	if req.PreferredEducationLevels != nil && len(*req.PreferredEducationLevels) > 0 {
		eduIn = *req.PreferredEducationLevels
	}
	educationLevelsTyped, err := validation.ParsePreferredEducationLevels(eduIn)
	if err != nil {
		return appError.ErrInvalidEducationLevel
	}
	preferences.PreferredEducationLevels = pq.StringArray(validation.EducationLevelsToStrings(educationLevelsTyped))

	// Home districts
	districtsIn := []string{"any"}
	if req.PreferredHomeDistricts != nil && len(*req.PreferredHomeDistricts) > 0 {
		districtsIn = *req.PreferredHomeDistricts
	}
	homeDistrictsTyped, err := validation.ParsePreferredHomeDistricts(districtsIn)
	if err != nil {
		return appError.ErrInvalidHomeDistrict
	}
	preferences.PreferredHomeDistricts = pq.StringArray(validation.HomeDistrictsToStrings(homeDistrictsTyped))

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
			patch["min_age_years"] = int16(*req.MinAgeYears)
			patch["max_age_years"] = int16(*req.MaxAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	} else if req.MinAgeYears != nil {
		if validation.IsValidAge(int(*req.MinAgeYears)) && int(*req.MinAgeYears) <= int(existingPartnerPreferences.MaxAgeYears) {
			patch["min_age_years"] = int16(*req.MinAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	} else if req.MaxAgeYears != nil {
		if validation.IsValidAge(int(*req.MaxAgeYears)) && int(*req.MaxAgeYears) >= int(existingPartnerPreferences.MinAgeYears) {
			patch["max_age_years"] = int16(*req.MaxAgeYears)
		} else {
			return appError.ErrInvalidAgeRange
		}
	}

	// Handle height range updates
	if req.MinHeightCM != nil && req.MaxHeightCM != nil {
		if validation.IsValidHeightRange(int(*req.MinHeightCM), int(*req.MaxHeightCM)) {
			patch["min_height_cm"] = int16(*req.MinHeightCM)
			patch["max_height_cm"] = int16(*req.MaxHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	} else if req.MinHeightCM != nil {
		if validation.IsValidHeight(int(*req.MinHeightCM)) && int(*req.MinHeightCM) <= int(existingPartnerPreferences.MaxHeightCm) {
			patch["min_height_cm"] = int16(*req.MinHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	} else if req.MaxHeightCM != nil {
		if validation.IsValidHeight(int(*req.MaxHeightCM)) && int(*req.MaxHeightCM) >= int(existingPartnerPreferences.MinHeightCm) {
			patch["max_height_cm"] = int16(*req.MaxHeightCM)
		} else {
			return appError.ErrInvalidHeightRange
		}
	}

	// Handle physically challenged preference
	if req.AcceptPhysicallyChallenged != nil {
		patch["accept_physically_challenged"] = *req.AcceptPhysicallyChallenged
	}

	if req.PreferredCommunities != nil && len(*req.PreferredCommunities) > 0 {
		communitiesTyped, err := validation.ParsePreferredCommunities(*req.PreferredCommunities)
		if err != nil {
			return appError.ErrInvalidCommunity
		}
		patch["preferred_communities"] = pq.StringArray(validation.CommunitiesToStrings(communitiesTyped))
	}

	if  req.PreferredMaritalStatus != nil && len(*req.PreferredMaritalStatus) > 0 {
		maritalTyped, err := validation.ParsePreferredMaritalStatuses(*req.PreferredMaritalStatus)
		if err != nil {
			return appError.ErrInvalidMaritalStatus
		}
		patch["preferred_marital_status"] = pq.StringArray(validation.MaritalStatusesToStrings(maritalTyped))
	}

	if req.PreferredProfessions != nil && len(*req.PreferredProfessions) > 0 {
		professionsTyped, err := validation.ParsePreferredProfessions(*req.PreferredProfessions)
    if err != nil {
        return appError.ErrInvalidProfession
    }
    patch["preferred_professions"] = pq.StringArray(validation.ProfessionsToStrings(professionsTyped))
	}

	if  req.PreferredProfessionTypes != nil && len(*req.PreferredProfessionTypes) > 0 {
		profTypesTyped, err := validation.ParsePreferredProfessionTypes(*req.PreferredProfessionTypes)
    if err != nil {
        return appError.ErrInvalidProfessionType
    }
    patch["preferred_profession_types"] = pq.StringArray(validation.ProfessionTypesToStrings(profTypesTyped))
	}

	if req.PreferredEducationLevels != nil && len(*req.PreferredEducationLevels) > 0 {
		eduTyped, err := validation.ParsePreferredEducationLevels(*req.PreferredEducationLevels)
    if err != nil {
        return appError.ErrInvalidEducationLevel
    }
    patch["preferred_education_levels"] = pq.StringArray(validation.EducationLevelsToStrings(eduTyped))
	}

	if req.PreferredHomeDistricts != nil && len(*req.PreferredHomeDistricts) > 0 {
		districtsTyped, err := validation.ParsePreferredHomeDistricts(*req.PreferredHomeDistricts)
    if err != nil {
        return appError.ErrInvalidHomeDistrict
    }
    patch["preferred_home_districts"] = pq.StringArray(validation.HomeDistrictsToStrings(districtsTyped))
	}

	if len(patch) > 0 {
		patch["updated_at"] = time.Now().UTC()
	}

	return u.partnerPreferencesRepository.PatchPartnerPreferences(ctx, userProfileID, patch)
}
