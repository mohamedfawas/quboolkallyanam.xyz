package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"gorm.io/gorm"
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

func (r *userProfileRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ? AND is_deleted = false", userID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appError.ErrUserProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) UpdateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error {
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userProfile.UserID).
		Updates(userProfile).Error
}

func (r *userProfileRepository) GetUserProfileByID(ctx context.Context, id uint) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("id = ? AND is_deleted = false", id).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appError.ErrUserProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}
