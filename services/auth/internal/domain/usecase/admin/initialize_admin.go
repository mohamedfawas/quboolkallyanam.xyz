package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) InitializeDefaultAdmin(ctx context.Context, defaultEmail, defaultPassword string) error {
	exists, err := u.adminRepository.CheckAdminExists(ctx, defaultEmail)
	if err != nil {
		return fmt.Errorf("failed to check if admin exists: %w", err)
	}
	if exists {
		return nil
	}

	if !validation.IsValidEmail(defaultEmail) {
		return apperrors.ErrInvalidEmail
	}

	if !validation.IsValidPassword(defaultPassword, validation.DefaultPasswordRequirements()) {
		return apperrors.ErrInvalidPassword
	}

	hasedPassword, err := hash.HashPassword(defaultPassword)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	admin := &entity.Admin{
		Email:        defaultEmail,
		PasswordHash: hasedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.adminRepository.CreateAdmin(ctx, admin); err != nil {
		return fmt.Errorf("failed to intialize admin account to database: %w", err)
	}

	return nil
}
