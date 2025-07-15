package pendingregistration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/otp"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

func (u *pendingRegistrationUsecase) generateAndStoreOTP(ctx context.Context, email string, config *config.Config) (string, error) {
	otp, err := otp.GenerateNumericOTP(6)
	if err != nil {
		log.Printf("failed to generate otp: %v", err)
		return "", fmt.Errorf("failed to generate otp: %w", err)
	}

	key := fmt.Sprintf("%s%s", constants.RedisPrefixOTP, email)
	if err := u.otpRepository.StoreOTP(ctx, key, otp, time.Minute*time.Duration(config.Auth.OTPExpiryMinutes)); err != nil {
		log.Printf("failed to store otp: %v", err)
		return "", fmt.Errorf("failed to store otp: %w", err)
	}

	return otp, nil
}

func (u *pendingRegistrationUsecase) validateOTP(ctx context.Context, inputOTP, OTPKey string) (bool, error) {
	storedOTP, err := u.otpRepository.GetOTP(ctx, OTPKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, appErrors.ErrOTPNotFound
		}
		log.Printf("failed to retrieve otp: %v", err)
		return false, fmt.Errorf("failed to retrieve otp: %w", err)
	}

	if inputOTP != storedOTP {
		log.Printf("invalid otp: %s", inputOTP)
		return false, nil
	}
	return true, nil
}
