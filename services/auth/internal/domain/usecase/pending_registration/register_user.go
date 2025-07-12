package pendingregistration

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp"
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

	logger.Log.Info("RegisterUser request received in pendingRegistrationUsecase", "email ", req.Email, "phone ", req.Phone)
	logger.Log.Info("Checking if email is already registered")
	existsEmail, err := u.userRepository.IsRegistered(ctx, "email", req.Email)
	if err != nil {
		logger.Log.Error("Error checking email", "error", err)
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existsEmail {
		logger.Log.Info("Email already registered", "email ", req.Email)
		return appErrors.ErrEmailAlreadyExists
	}

	logger.Log.Info("Checking if phone is already registered")
	existsPhone, err := u.userRepository.IsRegistered(ctx, "phone", req.Phone)
	if err != nil {
		logger.Log.Error("Error checking phone", "error", err)
		return fmt.Errorf("failed to check phone: %w", err)
	}
	if existsPhone {
		logger.Log.Info("Phone already registered", "phone ", req.Phone)
		return appErrors.ErrPhoneAlreadyExists
	}

	logger.Log.Info("Checking if email is already pending registration")
	pendingEmail, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Error checking pending registration by email", "error", err)
		return fmt.Errorf("failed to get existing pending registration: %w", err)
	}
	if pendingEmail != nil {
		err := u.pendingRegistrationRepository.DeletePendingRegistration(ctx, pendingEmail.ID)
		if err != nil {
			logger.Log.Error("Error deleting existing pending registration by email", "error", err)
			return fmt.Errorf("failed to delete existing pending registration: %w", err)
		}
	}

	pendingPhone, err := u.pendingRegistrationRepository.GetPendingRegistration(ctx, "phone", req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Error checking pending registration by phone", "error", err)
		return fmt.Errorf("failed to get existing pending registration: %w", err)
	}
	if pendingPhone != nil {
		err := u.pendingRegistrationRepository.DeletePendingRegistration(ctx, pendingPhone.ID)
		if err != nil {
			logger.Log.Error("Error deleting existing pending registration by phone", "error", err)
			return fmt.Errorf("failed to delete existing pending registration: %w", err)
		}
	}

	// Using Bcrypt default cost of 12
	logger.Log.Info("Hashing password in pendingRegistrationUsecase")
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("Error hashing password", "error", err)
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

	logger.Log.Info("Pending registration created in pendingRegistrationUsecase and sending to repository", "email ", req.Email, "phone ", req.Phone)
	if err := u.pendingRegistrationRepository.CreatePendingRegistration(ctx, pendingRegistration); err != nil {
		logger.Log.Error("Error creating pending registration", "error", err)
		return fmt.Errorf("failed to create pending registration: %w", err)
	}

	logger.Log.Info("Generating and storing OTP in pendingRegistrationUsecase", "email ", req.Email, "phone ", req.Phone)
	otp, err := u.generateAndStoreOTP(ctx, req.Email, config)
	if err != nil {
		logger.Log.Error("Error generating and storing OTP", "error", err)
		return err
	}

	logger.Log.Info("Sending OTP verification email in pendingRegistrationUsecase", "email ", req.Email, "phone ", req.Phone)
	if err := u.smtpClient.SendEmailByType(smtp.EmailRequest{
		To:      req.Email,
		Type:    smtp.EmailTypeOTPVerification,
		Subject: "Qubool Kallyanam - User Registration OTP Verification",
		Payload: map[string]string{
			"email":         req.Email,
			"otp":           otp,
			"expiryMinutes": strconv.Itoa(config.Auth.OTPExpiryMinutes),
		},
	}); err != nil {
		logger.Log.Error("Error sending OTP verification email", "error", err)
		return fmt.Errorf("failed to send OTP verification email: %w", err)
	}

	logger.Log.Info("Successfully sent OTP verification email in pendingRegistrationUsecase", "email ", req.Email, "phone ", req.Phone)
	return nil
}
