package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *postgres.Client
}

func NewUserRepository(db *postgres.Client) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context,
	field, value string) (*entity.User, error) {

	var user entity.User
	allowed := map[string]bool{"email": true, "phone": true, "id": true}
	if !allowed[field] {
		return nil, fmt.Errorf("invalid field %q for GetUser", field)
	}

	err := r.db.GormDB.
		WithContext(ctx).
		Where(map[string]interface{}{field: value}).
		First(&user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context,
	user *entity.User) error {

	return r.db.GormDB.
		WithContext(ctx).
		Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context,
	user *entity.User) error {

	return r.db.GormDB.
		WithContext(ctx).
		Save(user).Error
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, userID string) error {
	return r.db.GormDB.
		WithContext(ctx).
		Delete(&entity.User{}, "id = ?", userID).
		Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context,
	userID string,
	now time.Time) error {

	return r.db.GormDB.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"updated_at":    now,
		}).Error
}

func (r *userRepository) UpdatePremiumUntil(ctx context.Context,
	userID string,
	premiumUntil time.Time,
	now time.Time) error {

	return r.db.GormDB.
		WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"premium_until": premiumUntil,
			"updated_at":    now,
		}).Error
}

func (r *userRepository) IsRegistered(ctx context.Context,
	field, value string) (bool, error) {

	allowed := map[string]bool{"email": true, "phone": true}
	if !allowed[field] {
		return false, fmt.Errorf("invalid field %q for IsRegistered", field)
	}

	var count int64
	if err := r.db.GormDB.WithContext(ctx).
		Model(&entity.User{}).
		Where(map[string]interface{}{field: value}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) BlockUser(ctx context.Context,
	field, value string) error {
	
	allowed := map[string]bool{"email": true, "phone": true, "id": true}
	if !allowed[field] {
		return fmt.Errorf("invalid field %q for BlockUser", field)
	}

	return r.db.GormDB.WithContext(ctx).
		Model(&entity.User{}).
		Where(map[string]interface{}{field: value}).
		Updates(map[string]interface{}{"is_blocked": true, "updated_at": time.Now().UTC()}).Error
}


func (r *userRepository) GetUsers(
	ctx context.Context, 
	page, limit int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.GormDB.WithContext(ctx).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&users).Error
	return users, err
}