package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserLogin(
	ctx context.Context,
	req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {

	if !validation.IsValidEmail(req.Email) {
		return nil, apperrors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, apperrors.ErrInvalidPassword
	}

	return u.authClient.UserLogin(ctx, req)
}
