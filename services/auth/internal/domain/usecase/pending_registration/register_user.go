package pendingregistration

import (
	"context"
	"errors"
	"fmt"
	"time"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *pendingRegistrationUsecase) RegisterUser(ctx context.Context,
	req *entity.UserRegistrationRequest,
	config *config.Config,
) error {
	existsEmail, err := u.userRepository.IsRegistered(ctx, "email", req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existsEmail {
		return appErrors.ErrEmailAlreadyExists
	}

	existsPhone, err := u.userRepository.IsRegistered(ctx, "phone", req.Phone)
	if err != nil {
		return fmt.Errorf("failed to check phone: %w", err)
	}
	if existsPhone {
		return appErrors.ErrPhoneAlreadyExists
	}

	pendingEmail, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get existing pending registration: %w", err)
	}
	if pendingEmail != nil {
		err := u.pendingRegistrationRepository.DeletePendingRegistration(ctx, pendingEmail.ID)
		if err != nil {
			return fmt.Errorf("failed to delete existing pending registration: %w", err)
		}
	}

	pendingPhone, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "phone", req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get existing pending registration: %w", err)
	}
	if pendingPhone != nil {
		err := u.pendingRegistrationRepository.DeletePendingRegistration(ctx, pendingPhone.ID)
		if err != nil {
			return fmt.Errorf("failed to delete existing pending registration: %w", err)
		}
	}

	// Using Bcrypt default cost of 12
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := timeutil.NowIST()
	pendingRegistration := &entity.PendingRegistration{
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		ExpiresAt:    now.Add(time.Hour * time.Duration(config.Auth.PendingRegistrationExpiryHours)),
	}

	if err := u.pendingRegistrationRepository.CreatePendingRegistration(ctx, pendingRegistration); err != nil {
		return fmt.Errorf("failed to create pending registration: %w", err)
	}

	otp, err := u.generateAndStoreOTP(ctx, req.Email, config)
	if err != nil {
		return err
	}

	otpEvent := authevents.UserOTPRequestedEvent{
		Email:         req.Email,
		OTP:           otp,
		ExpiryMinutes: config.Auth.OTPExpiryMinutes,
	}

	if err := u.eventPublisher.PublishUserOTPRequested(ctx, otpEvent); err != nil {
		return fmt.Errorf("failed to publish OTP requested event: %w", err)
	}

	return nil
}
