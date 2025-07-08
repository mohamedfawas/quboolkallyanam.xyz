package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/vaidation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *AuthUsecase) UserDelete(ctx context.Context, req dto.UserDeleteRequest) error {
	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return errors.ErrInvalidPassword
	}

	return u.authClient.UserDelete(ctx, req)
}
