package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) UpdateUserPartnerPreferences(
	ctx context.Context,
	operationType string,
	req dto.PartnerPreferencePatchRequest) error {

	updatePartnerPreferenceRequest := dto.UpdatePartnerPreferenceRequest{}

	if operationType != constants.CreateOperationType && operationType != constants.UpdateOperationType {
		return apperrors.ErrInvalidOperationType
	}

	if req.MinAgeYears != nil {
		if !validation.IsValidAge(pointerutil.GetIntValue(req.MinAgeYears)) {
			return apperrors.ErrInvalidAgeRange
		}
		updatePartnerPreferenceRequest.MinAgeYears = req.MinAgeYears
	}

	if req.MaxAgeYears != nil {
		if !validation.IsValidAge(pointerutil.GetIntValue(req.MaxAgeYears)) {
			return apperrors.ErrInvalidAgeRange
		}
		updatePartnerPreferenceRequest.MaxAgeYears = req.MaxAgeYears
	}

	if req.MinAgeYears != nil && req.MaxAgeYears != nil {
		if !validation.IsValidAgeRange(pointerutil.GetIntValue(req.MinAgeYears), pointerutil.GetIntValue(req.MaxAgeYears)) {
			return apperrors.ErrInvalidAgeRange
		}
		updatePartnerPreferenceRequest.MinAgeYears = req.MinAgeYears
		updatePartnerPreferenceRequest.MaxAgeYears = req.MaxAgeYears
	}

	if req.MinHeightCM != nil {
		if !validation.IsValidHeight(pointerutil.GetIntValue(req.MinHeightCM)) {
			return apperrors.ErrInvalidHeightRange
		}
		updatePartnerPreferenceRequest.MinHeightCM = req.MinHeightCM
	}

	if req.MaxHeightCM != nil {
		if !validation.IsValidHeight(pointerutil.GetIntValue(req.MaxHeightCM)) {
			return apperrors.ErrInvalidHeightRange
		}
		updatePartnerPreferenceRequest.MaxHeightCM = req.MaxHeightCM
	}

	if req.MinHeightCM != nil && req.MaxHeightCM != nil {
		if !validation.IsValidHeightRange(pointerutil.GetIntValue(req.MinHeightCM), pointerutil.GetIntValue(req.MaxHeightCM)) {
			return apperrors.ErrInvalidHeightRange
		}
		updatePartnerPreferenceRequest.MinHeightCM = req.MinHeightCM
		updatePartnerPreferenceRequest.MaxHeightCM = req.MaxHeightCM

	}

	if req.PreferredCommunities != nil {
		communities, err := validation.ParsePreferredCommunities(*req.PreferredCommunities)
		if err != nil {
			return apperrors.ErrInvalidCommunity
		}
		updatePartnerPreferenceRequest.PreferredCommunities = &communities
	}

	if req.PreferredEducationLevels != nil {
		educationLevels, err := validation.ParsePreferredEducationLevels(*req.PreferredEducationLevels)
		if err != nil {
			return apperrors.ErrInvalidEducationLevel
		}
		updatePartnerPreferenceRequest.PreferredEducationLevels = &educationLevels
	}

	if req.PreferredHomeDistricts != nil {
		homeDistricts, err := validation.ParsePreferredHomeDistricts(*req.PreferredHomeDistricts)
		if err != nil {
			return apperrors.ErrInvalidHomeDistrict
		}
		updatePartnerPreferenceRequest.PreferredHomeDistricts = &homeDistricts
	}

	if req.PreferredMaritalStatus != nil {
		maritalStatuses, err := validation.ParsePreferredMaritalStatuses(*req.PreferredMaritalStatus)
		if err != nil {
			return apperrors.ErrInvalidMaritalStatus
		}
		updatePartnerPreferenceRequest.PreferredMaritalStatus = &maritalStatuses
	}

	if req.PreferredProfessions != nil {
		professions, err := validation.ParsePreferredProfessions(*req.PreferredProfessions)
		if err != nil {
			return apperrors.ErrInvalidProfession
		}
		updatePartnerPreferenceRequest.PreferredProfessions = &professions
	}

	if req.PreferredProfessionTypes != nil {
		professionTypes, err := validation.ParsePreferredProfessionTypes(*req.PreferredProfessionTypes)
		if err != nil {
			return apperrors.ErrInvalidProfessionType
		}
		updatePartnerPreferenceRequest.PreferredProfessionTypes = &professionTypes
	}

	return u.userClient.UpdateUserPartnerPreferences(ctx, operationType, updatePartnerPreferenceRequest)
}
