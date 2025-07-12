package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"

	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/vaidation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	if !validation.IsValidEmail(req.Email) {
		return nil, errors.ErrInvalidEmail
	}

	if !validation.IsValidPhoneNumber(req.Phone) {
		return nil, errors.ErrInvalidPhoneNumber
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, errors.ErrInvalidPassword
	}

	return u.authClient.UserRegister(ctx, req)
}
