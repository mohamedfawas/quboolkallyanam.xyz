package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) DeleteProfilePhoto(ctx context.Context, req dto.DeleteProfilePhotoRequest) error {
	_, err := u.userClient.DeleteProfilePhoto(ctx, req)
	if err != nil {
		return err
	}
	return nil
}