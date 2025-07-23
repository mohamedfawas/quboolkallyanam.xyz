package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserVerification(
	ctx context.Context,
	req dto.UserVerificationRequest,
	config config.Config) (*dto.UserVerificationResponse, error) {

	if !validation.IsValidEmail(req.Email) {
		return nil, errors.ErrInvalidEmail
	}

	if !validation.IsValidOTP(req.OTP, config.OTP.Length) {
		return nil, errors.ErrInvalidOTP
	}

	return u.authClient.UserVerification(ctx, req)
}
