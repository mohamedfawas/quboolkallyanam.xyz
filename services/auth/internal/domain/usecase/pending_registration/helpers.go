package pendingregistration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/otp"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

func (u *pendingRegistrationUsecase) generateAndStoreOTP(ctx context.Context, email string, config *config.Config) (string, error) {
	otp, err := otp.GenerateNumericOTP(constants.DefaultOTPLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate otp: %w", err)
	}

	key := fmt.Sprintf("%s%s", constants.RedisPrefixOTP, email)
	if err := u.otpRepository.StoreOTP(ctx, key, otp, time.Minute*time.Duration(config.Auth.OTPExpiryMinutes)); err != nil {
		return "", fmt.Errorf("failed to store otp: %w", err)
	}

	return otp, nil
}

func (u *pendingRegistrationUsecase) validateOTP(ctx context.Context, inputOTP, OTPKey string) (bool, error) {
	storedOTP, err := u.otpRepository.GetOTP(ctx, OTPKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, apperrors.ErrOTPNotFound
		}
		return false, fmt.Errorf("failed to retrieve otp: %w", err)
	}

	if inputOTP != storedOTP {
		return false, nil
	}
	return true, nil
}
