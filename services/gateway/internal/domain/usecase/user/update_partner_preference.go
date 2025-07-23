package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) UpdateUserPartnerPreferences(
	ctx context.Context,
	operationType string,
	req dto.PartnerPreferencePatchRequest) error {

	if operationType != constants.CreateOperationType && operationType != constants.UpdateOperationType {
		return appError.ErrInvalidOperationType
	}

	if req.MinAgeYears != nil {
		if !validation.IsValidAge(req.MinAgeYears) {
			return appError.ErrInvalidAgeRange
		}
	}

	if req.MaxAgeYears != nil {
		if !validation.IsValidAge(req.MaxAgeYears) {
			return appError.ErrInvalidAgeRange
		}
	}

	if req.MinAgeYears != nil && req.MaxAgeYears != nil {
		if !validation.IsValidAgeRange(req.MinAgeYears, req.MaxAgeYears) {
			return appError.ErrInvalidAgeRange
		}
	}

	if req.MinHeightCM != nil {
		if !validation.IsValidHeight(req.MinHeightCM) {
			return appError.ErrInvalidHeightRange
		}
	}

	if req.MaxHeightCM != nil {
		if !validation.IsValidHeight(req.MaxHeightCM) {
			return appError.ErrInvalidHeightRange
		}
	}

	if req.MinHeightCM != nil && req.MaxHeightCM != nil {
		if !validation.IsValidHeightRange(req.MinHeightCM, req.MaxHeightCM) {
			return appError.ErrInvalidHeightRange
		}
	}

	if req.PreferredCommunities != nil {
		for _, community := range *req.PreferredCommunities {
			if !validation.IsValidCommunity(community) {
				return appError.ErrInvalidCommunity
			}
		}
	}

	if req.PreferredEducationLevels != nil {
		for _, educationLevel := range *req.PreferredEducationLevels {
			if !validation.IsValidEducationLevel(educationLevel) {
				return appError.ErrInvalidEducationLevel
			}
		}
	}

	if req.PreferredHomeDistricts != nil {
		for _, homeDistrict := range *req.PreferredHomeDistricts {
			if !validation.IsValidHomeDistrict(homeDistrict) {
				return appError.ErrInvalidHomeDistrict
			}
		}
	}

	if req.PreferredMaritalStatus != nil {
		for _, maritalStatus := range *req.PreferredMaritalStatus {
			if !validation.IsValidMaritalStatus(maritalStatus) {
				return appError.ErrInvalidMaritalStatus
			}
		}
	}

	if req.PreferredProfessions != nil {
		for _, profession := range *req.PreferredProfessions {
			if !validation.IsValidProfession(profession) {
				return appError.ErrInvalidProfession
			}
		}
	}

	if req.PreferredProfessionTypes != nil {
		for _, professionType := range *req.PreferredProfessionTypes {
			if !validation.IsValidProfessionType(professionType) {
				return appError.ErrInvalidProfessionType
			}
		}
	}

	return u.userClient.UpdateUserPartnerPreferences(ctx, operationType, req)
}
