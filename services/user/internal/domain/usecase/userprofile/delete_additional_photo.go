package userprofile

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
)

func (u *userProfileUsecase) DeleteAdditionalPhoto(
	ctx context.Context,
	userID uuid.UUID,
	displayOrder int32) error {

	userImage, err := u.userImageRepository.GetUserImage(ctx, userID, displayOrder)
	if err != nil {
		return err
	}
	if userImage == nil {
		return apperrors.ErrDisplayOrderNotOccupied
	}

	err = u.photoStorage.DeletePhoto(ctx, userImage.ObjectKey)
	if err != nil {
		return err
	}

	err = u.userImageRepository.DeleteUserImage(ctx, userID, displayOrder)
	if err != nil {
		return err
	}

	return nil
}
