package admin

import (
	"context"
	"fmt"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/vaidation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) InitializeDefaultAdmin(ctx context.Context, defaultEmail, defaultPassword string) error {
	exists, err := u.adminRepository.CheckAdminExists(ctx, defaultEmail)
	if err != nil {
		return fmt.Errorf("failed to check if admin exists: %w", err)
	}
	if exists {
		logger.Log.Info("Admin already exists, skipping default admin creation")
		return nil
	}

	if validation.IsValidEmail(defaultEmail) {
		logger.Log.Error("Provided email for admin initialization is invalid", defaultEmail)
		return errors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(defaultPassword, validation.DefaultPasswordRequirements()) {
		logger.Log.Error("Provided password for admin initialization is not strong enough", defaultPassword)
		return errors.ErrInvalidPassword
	}

	hasedPassword, err := hash.HashPassword(defaultPassword)
	if err != nil {
		logger.Log.Error("Failed to hash given admin password", err)
		return errors.ErrHashGenerationFailed
	}

	now := timeutil.NowIST()
	admin := &entity.Admin{
		Email:        defaultEmail,
		PasswordHash: hasedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.adminRepository.CreateAdmin(ctx, admin); err != nil {
		return fmt.Errorf("failed to intialize admin account to database: %w", err)
	}

	logger.Log.Info("Default admin account initialized successfully", "email", defaultEmail)

	return nil
}
