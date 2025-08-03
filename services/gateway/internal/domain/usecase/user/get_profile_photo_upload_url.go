package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetProfilePhotoUploadURL(
	ctx context.Context,
	req dto.GetProfilePhotoUploadURLRequest) (*dto.GetProfilePhotoUploadURLResponse, error) {

	if !validation.IsValidImageType(req.ContentType) {
		return nil, apperrors.ErrInvalidImageType
	}

	return u.userClient.GetProfilePhotoUploadURL(ctx, req)
}
