package userprojection

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

func (u *userProjectionUsecase) CreateOrUpdateUserProjection(
	ctx context.Context,
	userProjection *entity.UserProjection) error {

	existingUserProjection, err := u.userProjectionRepository.GetUserProjectionByID(ctx, userProjection.UserProfileID)
	if err != nil {
		return err
	}

	if existingUserProjection != nil {
		if existingUserProjection.Email == userProjection.Email && existingUserProjection.FullName == userProjection.FullName {
			return nil
		}
		existingUserProjection.Email = userProjection.Email
		existingUserProjection.FullName = userProjection.FullName
		existingUserProjection.UpdatedAt = time.Now().UTC()
		return u.userProjectionRepository.UpdateUserProjection(ctx, userProjection.UserProfileID, existingUserProjection)
	}

	return u.userProjectionRepository.CreateUserProjection(ctx, userProjection)
}
