package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetUserPartnerPreferences(ctx context.Context) (*dto.GetUserPartnerPreferencesResponse, error) {
	return u.userClient.GetUserPartnerPreferences(ctx)
}