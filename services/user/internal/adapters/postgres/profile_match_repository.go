package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"gorm.io/gorm"
)

type profileMatchRepository struct {
	db *postgres.Client
}

func NewProfileMatchRepository(db *postgres.Client) repository.ProfileMatchRepository {
	return &profileMatchRepository{db: db}
}

func (r *profileMatchRepository) GetExistingMatch(ctx context.Context, userID uuid.UUID, targetID uuid.UUID) (*entity.ProfileMatch, error) {
	var match entity.ProfileMatch

	err := r.db.GormDB.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND is_deleted = false", userID, targetID).
		First(&match).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &match, nil
}

func (r *profileMatchRepository) CreateMatchAction(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error {
	match := entity.ProfileMatch{
		UserID:   userID,
		TargetID: targetID,
		IsLiked:  isLiked,
	}

	return r.db.GormDB.WithContext(ctx).Create(&match).Error
}

func (r *profileMatchRepository) CreateMatchActionTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error {
	match := entity.ProfileMatch{
		UserID:   userID,
		TargetID: targetID,
		IsLiked:  isLiked,
	}

	return tx.WithContext(ctx).Create(&match).Error
}

func (r *profileMatchRepository) UpdateMatchAction(
	ctx context.Context,
	userID uuid.UUID,
	targetID uuid.UUID,
	isLiked bool) error {

	now := time.Now().UTC()
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.ProfileMatch{}).
		Where("user_id = ? AND target_id = ? AND is_deleted = false",
			userID, targetID).
		Updates(map[string]interface{}{
			"is_liked":   isLiked,
			"updated_at": now,
		}).Error
}

func (r *profileMatchRepository) UpdateMatchActionTx(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	targetID uuid.UUID,
	isLiked bool) error {

	now := time.Now().UTC()
	return tx.WithContext(ctx).
		Model(&entity.ProfileMatch{}).
		Where("user_id = ? AND target_id = ? AND is_deleted = false",
			userID, targetID).
		Updates(map[string]interface{}{
			"is_liked":   isLiked,
			"updated_at": now,
		}).Error
}
