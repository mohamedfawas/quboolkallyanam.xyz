package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetAdditionalPhotos(ctx context.Context) (*dto.GetAdditionalPhotosResponse, error) {
	return u.userClient.GetAdditionalPhotos(ctx)
}