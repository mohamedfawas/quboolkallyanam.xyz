package postgres

import (
	"context"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
)

type userRepository struct {
	db *postgres.Client
}

func NewUserRepository(db *postgres.Client) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context, field, value string) (*entity.User, error) {
	var user entity.User
	if err := r.db.GormDB.WithContext(ctx).Where(field+" = ?", value).First(&user).Error; err != nil {
		log.Printf("GetUser error in user repository: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.GormDB.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.GormDB.WithContext(ctx).Save(user).Error
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, userID string, now time.Time) error {
	return r.db.GormDB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": now,
		"deleted_at": now,
	}).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID string, now time.Time) error {
	return r.db.GormDB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login_at": now,
		"updated_at":    now,
	}).Error
}

func (r *userRepository) UpdatePremiumUntil(ctx context.Context, userID string, premiumUntil time.Time, now time.Time) error {
	return r.db.GormDB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"premium_until": premiumUntil,
		"updated_at":    now,
	}).Error
}

func (r *userRepository) IsRegistered(ctx context.Context, field, value string) (bool, error) {
	var count int64
	if err := r.db.GormDB.WithContext(ctx).Model(&entity.User{}).Where(field+" = ?", value).Count(&count).Error; err != nil {
		log.Printf("IsRegistered error in user repository: %v", err)
		return false, err
	}
	return count > 0, nil
}
