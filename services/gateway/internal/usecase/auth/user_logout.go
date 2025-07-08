package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *AuthUsecase) UserLogout(ctx context.Context, req dto.UserLogoutRequest) error {
	return u.authClient.UserLogout(ctx, req)
}
