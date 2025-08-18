package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error {

	if req.FullName != nil {
		if !validation.IsValidFullName(*req.FullName) {
			return apperrors.ErrInvalidFullName
		}
	}

	if req.Community != nil {
		if !validation.IsValidCommunity(*req.Community) {
			return apperrors.ErrInvalidCommunity
		}
	}

	if req.MaritalStatus != nil {
		if !validation.IsValidMaritalStatus(*req.MaritalStatus) {
			return apperrors.ErrInvalidMaritalStatus
		}
	}

	if req.Profession != nil {
		if !validation.IsValidProfession(*req.Profession) {
			return apperrors.ErrInvalidProfession
		}
	}

	if req.ProfessionType != nil {
		if !validation.IsValidProfessionType(*req.ProfessionType) {
			return apperrors.ErrInvalidProfessionType
		}
	}

	if req.HighestEducationLevel != nil {
		if !validation.IsValidEducationLevel(*req.HighestEducationLevel) {
			return apperrors.ErrInvalidEducationLevel
		}
	}

	if req.HomeDistrict != nil {
		if !validation.IsValidHomeDistrict(*req.HomeDistrict) {
			return apperrors.ErrInvalidHomeDistrict
		}
	}

	if req.DateOfBirth != nil {
		if !validation.IsValidDateOfBirth(*req.DateOfBirth) {
			return apperrors.ErrInvalidDateOfBirth
		}
	}

	if req.HeightCm != nil {
		if !validation.IsValidHumanHeight(*req.HeightCm) {
			return apperrors.ErrInvalidHeight
		}
	}

	return u.userClient.UpdateUserProfile(ctx, req)
}
