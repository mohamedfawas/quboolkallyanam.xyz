package userprofile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/gcsutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) ConfirmAdditionalPhotoUpload(ctx context.Context,
	userID uuid.UUID,
	objectKey string,
	fileSize uint64) (string, error) {

	if !validation.IsValidImageFileSize(fileSize) {
		return "", apperrors.ErrInvalidImageFileSize
	}

	profile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}
	if profile == nil {
		return "", apperrors.ErrUserNotFound
	}

	displayOrder, err := gcsutil.ExtractDisplayOrder(objectKey)
	if err != nil {
		return "", fmt.Errorf("failed to extract display order: %w", err)
	}

	err = u.userImageRepository.CreateUserImage(ctx, &entity.UserImage{
		UserID:       userID,
		ObjectKey:    objectKey,
		DisplayOrder: int16(displayOrder),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user image: %w", err)
	}

	downloadURL, err := u.photoStorage.GetDownloadURL(ctx, objectKey, u.config.MediaStorage.URLExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}
	return downloadURL, nil
}
