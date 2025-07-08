package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/vaidation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *AuthUsecase) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	if !validation.IsValidEmail(req.Email) {
		return nil, errors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, errors.ErrInvalidPassword
	}

	return u.authClient.AdminLogin(ctx, req)
}
