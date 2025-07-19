package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
)

type userProfileRepository struct {
	db *postgres.Client
}

func NewUserProfileRepository(db *postgres.Client) repository.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) CreateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error {
	return r.db.GormDB.WithContext(ctx).
		Create(userProfile).Error
}

func (r *userProfileRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLogin time.Time) error {
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userID).
		Update("last_login", lastLogin).Error
}

func (r *userProfileRepository) ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var matched int64
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Count(&matched).Error
	if err != nil {
		return false, err
	}
	return matched > 0, nil
}
