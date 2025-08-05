package userprofile

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
)

func (u *userProfileUsecase) DeleteProfilePhoto(ctx context.Context, userID uuid.UUID) error {
	profile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return apperrors.ErrUserNotFound
	}

	if profile.ProfileImageKey == "" {
		return apperrors.ErrProfilePhotoNotFound
	}

	if err := u.photoStorage.DeletePhoto(ctx, profile.ProfileImageKey); err != nil {
		return err
	}

	profile.ProfileImageKey = ""
	if err := u.userProfileRepository.UpdateUserProfile(ctx, profile); err != nil {
		return err
	}

	return nil
}