package admin

import (
	"context"
	"fmt"
	"log"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
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
		log.Println("Admin already exists, skipping default admin creation")
		return nil
	}

	if !validation.IsValidEmail(defaultEmail) {
		log.Printf("Provided email for admin initialization is invalid: %v", defaultEmail)
		return appErrors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(defaultPassword, validation.DefaultPasswordRequirements()) {
		log.Printf("Provided password for admin initialization is not strong enough: %v", defaultPassword)
		return appErrors.ErrInvalidPassword
	}

	hasedPassword, err := hash.HashPassword(defaultPassword)
	if err != nil {
		log.Printf("Failed to hash given admin password: %v", err)
		return appErrors.ErrHashGenerationFailed
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

	log.Printf("Default admin account initialized successfully: %v", defaultEmail)

	return nil
}
