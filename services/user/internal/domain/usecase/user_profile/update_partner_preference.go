package userprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
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
		log.Printf("failed to get user profile: %v", err)
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	existingPartnerPreferences, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx,
		uint(existingProfile.ID))
	if err != nil {
		if err != appError.ErrPartnerPreferencesNotFound {
			log.Printf("failed to get partner preferences: %v", err)
			return fmt.Errorf("failed to get partner preferences: %w", err)
		}
	}

	if operationType == constants.CreateOperationType {
		if existingPartnerPreferences != nil {
			return appError.ErrPartnerPreferencesAlreadyExists
		}
		return u.createPartnerPreferences(ctx, existingProfile.ID, req)
	}

	if err == appError.ErrPartnerPreferencesNotFound {
		return appError.ErrPartnerPreferencesNotFound
	}

	return u.patchPartnerPreferences(ctx, uint(existingProfile.ID), req)
}

func (u *userProfileUsecase) createPartnerPreferences(
	ctx context.Context,
	userProfileID int64,
	req entity.UpdateUserPartnerPreferencesRequest) error {

	preferences := &entity.PartnerPreference{
		UserProfileID:              userProfileID,
		AcceptPhysicallyChallenged: true, // default value
	}

	if req.MinAgeYears != nil {
		preferences.MinAgeYears = req.MinAgeYears
	}
	if req.MaxAgeYears != nil {
		preferences.MaxAgeYears = req.MaxAgeYears
	}
	if req.MinHeightCM != nil {
		preferences.MinHeightCm = req.MinHeightCM
	}
	if req.MaxHeightCM != nil {
		preferences.MaxHeightCm = req.MaxHeightCM
	}
	if req.AcceptPhysicallyChallenged != nil {
		preferences.AcceptPhysicallyChallenged = *req.AcceptPhysicallyChallenged
	}

	if req.PreferredCommunities != nil {
		preferences.PreferredCommunities = convertToCommunityEnums(*req.PreferredCommunities)
	}
	if req.PreferredMaritalStatus != nil {
		preferences.PreferredMaritalStatus = convertToMaritalStatusEnums(*req.PreferredMaritalStatus)
	}
	if req.PreferredProfessions != nil {
		preferences.PreferredProfessions = convertToProfessionEnums(*req.PreferredProfessions)
	}
	if req.PreferredProfessionTypes != nil {
		preferences.PreferredProfessionTypes = convertToProfessionTypeEnums(*req.PreferredProfessionTypes)
	}
	if req.PreferredEducationLevels != nil {
		preferences.PreferredEducationLevels = convertToEducationLevelEnums(*req.PreferredEducationLevels)
	}
	if req.PreferredHomeDistricts != nil {
		preferences.PreferredHomeDistricts = convertToHomeDistrictEnums(*req.PreferredHomeDistricts)
	}

	return u.partnerPreferencesRepository.CreatePartnerPreferences(ctx, preferences)
}

func (u *userProfileUsecase) patchPartnerPreferences(
	ctx context.Context,
	userProfileID uint,
	req entity.UpdateUserPartnerPreferencesRequest) error {

	patch := make(map[string]interface{})

	if req.MinAgeYears != nil {
		patch["min_age_years"] = *req.MinAgeYears
	}
	if req.MaxAgeYears != nil {
		patch["max_age_years"] = *req.MaxAgeYears
	}
	if req.MinHeightCM != nil {
		patch["min_height_cm"] = *req.MinHeightCM
	}
	if req.MaxHeightCM != nil {
		patch["max_height_cm"] = *req.MaxHeightCM
	}
	if req.AcceptPhysicallyChallenged != nil {
		patch["accept_physically_challenged"] = *req.AcceptPhysicallyChallenged
	}

	if req.PreferredCommunities != nil {
		// Convert enum array to JSON string for proper JSONB storage
		communities := convertToCommunityEnums(*req.PreferredCommunities)
		jsonData, err := json.Marshal(communities)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred communities: %w", err)
		}
		patch["preferred_communities"] = string(jsonData)
	}
	if req.PreferredMaritalStatus != nil {
		maritalStatus := convertToMaritalStatusEnums(*req.PreferredMaritalStatus)
		jsonData, err := json.Marshal(maritalStatus)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred marital status: %w", err)
		}
		patch["preferred_marital_status"] = string(jsonData)
	}
	if req.PreferredProfessions != nil {
		professions := convertToProfessionEnums(*req.PreferredProfessions)
		jsonData, err := json.Marshal(professions)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred professions: %w", err)
		}
		patch["preferred_professions"] = string(jsonData)
	}
	if req.PreferredProfessionTypes != nil {
		professionTypes := convertToProfessionTypeEnums(*req.PreferredProfessionTypes)
		jsonData, err := json.Marshal(professionTypes)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred profession types: %w", err)
		}
		patch["preferred_profession_types"] = string(jsonData)
	}
	if req.PreferredEducationLevels != nil {
		educationLevels := convertToEducationLevelEnums(*req.PreferredEducationLevels)
		jsonData, err := json.Marshal(educationLevels)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred education levels: %w", err)
		}
		patch["preferred_education_levels"] = string(jsonData)
	}
	if req.PreferredHomeDistricts != nil {
		homeDistricts := convertToHomeDistrictEnums(*req.PreferredHomeDistricts)
		jsonData, err := json.Marshal(homeDistricts)
		if err != nil {
			return fmt.Errorf("failed to marshal preferred home districts: %w", err)
		}
		patch["preferred_home_districts"] = string(jsonData)
	}

	return u.partnerPreferencesRepository.PatchPartnerPreferences(ctx, userProfileID, patch)
}

func convertToCommunityEnums(strs []string) []entity.CommunityEnum {
	enums := make([]entity.CommunityEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.CommunityEnum(str)
	}
	return enums
}

func convertToMaritalStatusEnums(strs []string) []entity.MaritalStatusEnum {
	enums := make([]entity.MaritalStatusEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.MaritalStatusEnum(str)
	}
	return enums
}

func convertToProfessionEnums(strs []string) []entity.ProfessionEnum {
	enums := make([]entity.ProfessionEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.ProfessionEnum(str)
	}
	return enums
}

func convertToProfessionTypeEnums(strs []string) []entity.ProfessionTypeEnum {
	enums := make([]entity.ProfessionTypeEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.ProfessionTypeEnum(str)
	}
	return enums
}

func convertToEducationLevelEnums(strs []string) []entity.EducationLevelEnum {
	enums := make([]entity.EducationLevelEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.EducationLevelEnum(str)
	}
	return enums
}

func convertToHomeDistrictEnums(strs []string) []entity.HomeDistrictEnum {
	enums := make([]entity.HomeDistrictEnum, len(strs))
	for i, str := range strs {
		enums[i] = entity.HomeDistrictEnum(str)
	}
	return enums
}
