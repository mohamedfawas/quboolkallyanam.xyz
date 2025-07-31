package user

import (
	"context"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error {

	if req.FullName != nil {
		if !validation.IsValidFullName(*req.FullName) {
			return appError.ErrInvalidFullName
		}
	}

	if req.Community != nil {
		if !validation.IsValidCommunity(*req.Community) {
			return appError.ErrInvalidCommunity
		}
	}

	if req.MaritalStatus != nil {
		if !validation.IsValidMaritalStatus(*req.MaritalStatus) {
			return appError.ErrInvalidMaritalStatus
		}
	}

	if req.Profession != nil {
		if !validation.IsValidProfession(*req.Profession) {
			return appError.ErrInvalidProfession
		}
	}

	if req.ProfessionType != nil {
		if !validation.IsValidProfessionType(*req.ProfessionType) {
			return appError.ErrInvalidProfessionType
		}
	}

	if req.HighestEducationLevel != nil {
		if !validation.IsValidEducationLevel(*req.HighestEducationLevel) {
			return appError.ErrInvalidEducationLevel
		}
	}

	if req.HomeDistrict != nil {
		if !validation.IsValidHomeDistrict(*req.HomeDistrict) {
			return appError.ErrInvalidHomeDistrict
		}
	}

	if req.DateOfBirth != nil {
		if !validation.IsValidDateOfBirth(*req.DateOfBirth) {
			return appError.ErrInvalidDateOfBirth
		}
	}

	if req.HeightCm != nil {
		if !validation.IsValidHumanHeight(*req.HeightCm) {
			return appError.ErrInvalidHeight
		}
	}

	return u.userClient.UpdateUserProfile(ctx, req)
}
