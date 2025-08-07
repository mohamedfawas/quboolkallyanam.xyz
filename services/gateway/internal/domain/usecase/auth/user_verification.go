package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserVerification(
	ctx context.Context,
	req dto.UserVerificationRequest,
	config config.Config) (*dto.UserVerificationResponse, error) {

	if !validation.IsValidEmail(req.Email) {
		return nil, apperrors.ErrInvalidEmail
	}

	if !validation.IsValidOTP(req.OTP, config.OTP.Length) {
		return nil, apperrors.ErrInvalidOTP
	}

	return u.authClient.UserVerification(ctx, req)
}
