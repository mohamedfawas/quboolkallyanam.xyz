package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetAdditionalPhotoUploadURL(
	ctx context.Context,
	req dto.GetAdditionalPhotoUploadURLRequest) (*dto.GetAdditionalPhotoUploadURLResponse, error) {

	
	if !validation.IsValidImageType(req.ContentType) {
		return nil, apperrors.ErrInvalidImageType
	}

	if req.DisplayOrder < 1 || req.DisplayOrder > 3 {
		return nil, apperrors.ErrInvalidDisplayOrder
	}

	return u.userClient.GetAdditionalPhotoUploadURL(ctx, req)
}