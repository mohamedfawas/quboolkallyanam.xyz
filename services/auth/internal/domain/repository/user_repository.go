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
	UpdateLastLogin(ctx context.Context, userID string, now time.Time) error
	UpdatePremiumUntil(ctx context.Context, userID string, premiumUntil time.Time, now time.Time) error
	IsRegistered(ctx context.Context, field, value string) (bool, error)
	BlockUser(ctx context.Context, field, value string) error
	GetUsers(ctx context.Context, page, limit int) ([]*entity.User, error)
}
