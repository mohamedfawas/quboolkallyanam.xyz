package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *AuthUsecase) AdminLogout(ctx context.Context, req dto.AdminLogoutRequest) error {
	return u.authClient.AdminLogout(ctx, req)
}
