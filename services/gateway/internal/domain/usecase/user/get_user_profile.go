package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetUserProfile(ctx context.Context) (*dto.UserProfileRecommendation, error) {
	return u.userClient.GetUserProfile(ctx)
}