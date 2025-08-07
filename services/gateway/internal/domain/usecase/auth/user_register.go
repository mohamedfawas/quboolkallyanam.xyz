package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserRegister(
	ctx context.Context,
	req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {

	if !validation.IsValidEmail(req.Email) {
		return nil, apperrors.ErrInvalidEmail
	}

	if !validation.IsValidPhoneNumber(req.Phone) {
		return nil, apperrors.ErrInvalidPhoneNumber
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, apperrors.ErrInvalidPassword
	}

	response, err := u.authClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
