package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
)

type userImageRepository struct {
	db *postgres.Client
}

func NewUserImageRepository(db *postgres.Client) repository.UserImageRepository {
	return &userImageRepository{db: db}
}

func (r *userImageRepository) DisplayOrderOccupied(
	ctx context.Context,
	userID uuid.UUID,
	displayOrder int32) (bool, error) {
	var count int64
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserImage{}).
		Where("user_id = ? AND display_order = ?", userID, displayOrder).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userImageRepository) CreateUserImage(ctx context.Context, userImage *entity.UserImage) error {
	err := r.db.GormDB.WithContext(ctx).Create(userImage).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userImageRepository) DeleteUserImage(ctx context.Context, userID uuid.UUID, displayOrder int32) error {
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserImage{}).
		Where("user_id = ? AND display_order = ?", userID, displayOrder).
		Delete(&entity.UserImage{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userImageRepository) GetUserImage(ctx context.Context, userID uuid.UUID, displayOrder int32) (*entity.UserImage, error) {
	var userImage entity.UserImage
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserImage{}).
		Where("user_id = ? AND display_order = ?", userID, displayOrder).
		First(&userImage).Error
	if err != nil {
		return nil, err
	}
}