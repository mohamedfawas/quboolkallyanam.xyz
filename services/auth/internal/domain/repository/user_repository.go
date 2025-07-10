package repository

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type UserRepository interface {
	GetUser(ctx context.Context, field, value string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	SoftDeleteUser(ctx context.Context, userID string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	UpdatePremiumUntil(ctx context.Context, userID string, premiumUntil time.Time) error
	IsRegistered(ctx context.Context, field, value string) (bool, error)
	// GetUsers(ctx context.Context, params GetUsersParams) ([]*entity.User, int64, error)
}
