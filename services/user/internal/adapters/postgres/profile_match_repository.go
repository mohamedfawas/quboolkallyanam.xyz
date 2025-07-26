package postgres

import (
	"context"
	"log"
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
		log.Printf("failed to get existing match: %v", err)
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

func (r *profileMatchRepository) GetMatchedProfileIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var matchedIDs []uuid.UUID
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.ProfileMatch{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Pluck("target_id", &matchedIDs).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("failed to get matched profile IDs: %v", err)
		return nil, err
	}

	return matchedIDs, nil
}

func (r *profileMatchRepository) GetLikedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var likedIDs []uuid.UUID
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.ProfileMatch{}).
		Where("user_id = ? AND is_liked = true AND is_deleted = false", userID).
		Order("updated_at DESC").
		Pluck("target_id", &likedIDs).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("failed to get liked user IDs: %v", err)
		return nil, err
	}

	return likedIDs, nil
}

func (r *profileMatchRepository) GetPassedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var passedIDs []uuid.UUID
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.ProfileMatch{}).
		Where("user_id = ? AND is_liked = false AND is_deleted = false", userID).
		Order("updated_at DESC").
		Pluck("target_id", &passedIDs).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("failed to get passed user IDs: %v", err)
		return nil, err
	}

	return passedIDs, nil
}
