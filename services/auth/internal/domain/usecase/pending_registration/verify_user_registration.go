package pendingregistration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *pendingRegistrationUsecase) VerifyUserRegistration(ctx context.Context,
	email, otp string) error {

	pendingRegistration, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "email", email)
	if err != nil {
		return fmt.Errorf("failed to get pending registration: %w", err)
	}

	if pendingRegistration == nil {
		return apperrors.ErrPendingRegistrationNotFound
	}

	OTPKey := fmt.Sprintf("%s%s", constants.RedisPrefixOTP, email)
	valid, err := u.validateOTP(ctx, otp, OTPKey)
	if err != nil {
		if errors.Is(err, apperrors.ErrOTPNotFound) {
			return apperrors.ErrOTPNotFound
		}
		return fmt.Errorf("failed to validate otp: %w", err)
	}

	if !valid {
		return apperrors.ErrInvalidOTP
	}

	if err := u.otpRepository.DeleteOTP(ctx, OTPKey); err != nil {
		return fmt.Errorf("failed to delete after otp validation: %w", err)
	}

	now := time.Now().UTC()
	user := &entity.User{
		Email:         pendingRegistration.Email,
		Phone:         pendingRegistration.Phone,
		PasswordHash:  pendingRegistration.PasswordHash,
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := u.pendingRegistrationRepository.DeletePendingRegistration(ctx, pendingRegistration.ID); err != nil {
		return fmt.Errorf("failed to delete pending registration: %w", err)
	}

	return nil
}
