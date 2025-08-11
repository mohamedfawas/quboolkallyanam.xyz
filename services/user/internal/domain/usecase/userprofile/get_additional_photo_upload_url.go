package userprofile

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
)

func (u *userProfileUsecase) GetAdditionalPhotoUploadURL(ctx context.Context,
	userID uuid.UUID,
	displayOrder int32,
	contentType string) (*mediastorage.PhotoUploadURLResponse, error) {


	exists, err := u.userProfileRepository.ProfileExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	occupied, err := u.userImageRepository.DisplayOrderOccupied(ctx, userID, displayOrder)
	if err != nil {
		return nil, err
	}
	if occupied {
		return nil, apperrors.ErrImageDisplayOrderOccupied
	}

	response, err := u.photoStorage.GetAdditionalPhotoUploadURL(ctx,
		userID.String(),
		displayOrder,
		contentType,
		u.config.MediaStorage.URLExpiry)
	if err != nil {
		return nil, err
	}
	return response, nil
}
