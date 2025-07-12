package pendingregistration

import (
	"context"
	"errors"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	timeutil "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *pendingRegistrationUsecase) VerifyUserRegistration(ctx context.Context,
	email, otp string) error {

	pendingRegistration, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "email", email)
	if err != nil {
		return fmt.Errorf("failed to get pending registration: %w", err)
	}

	if pendingRegistration == nil {
		return appErrors.ErrPendingRegistrationNotFound
	}

	OTPKey := fmt.Sprintf("%s%s", constants.RedisPrefixOTP, email)
	valid, err := u.validateOTP(ctx, otp, OTPKey)
	if err != nil {
		if errors.Is(err, appErrors.ErrOTPNotFound) {
			return appErrors.ErrOTPNotFound
		}
		return fmt.Errorf("failed to validate otp: %w", err)
	}

	if !valid {
		return appErrors.ErrInvalidOTP
	}

	if err := u.otpRepository.DeleteOTP(ctx, OTPKey); err != nil {
		return fmt.Errorf("failed to delete after otp validation: %w", err)
	}

	now := timeutil.NowIST()
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
