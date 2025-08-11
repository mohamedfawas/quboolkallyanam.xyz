package userprofile

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
)

func (u *userProfileUsecase) GetProfilePhotoUploadURL(ctx context.Context, 
	userID uuid.UUID,
	 contentType string) (*mediastorage.PhotoUploadURLResponse, error) {
	if !validation.IsValidImageType(contentType) {
		return nil, apperrors.ErrInvalidImageType
	}

	exists, err := u.userProfileRepository.ProfileExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	response, err := u.photoStorage.GetProfilePhotoUploadURL(ctx, 
		userID.String(), 
		contentType, 
		u.config.MediaStorage.URLExpiry)
	if err != nil {
		return nil, err
	}
	return response, nil
}
