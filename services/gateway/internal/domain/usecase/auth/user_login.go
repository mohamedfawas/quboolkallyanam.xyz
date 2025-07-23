package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserLogin(
	ctx context.Context,
	req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {

	if !validation.IsValidEmail(req.Email) {
		return nil, errors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, errors.ErrInvalidPassword
	}

	return u.authClient.UserLogin(ctx, req)
}
