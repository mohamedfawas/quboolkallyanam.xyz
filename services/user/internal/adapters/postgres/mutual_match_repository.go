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

type mutualMatchRepository struct {
	db *postgres.Client
}

func NewMutualMatchRepository(db *postgres.Client) repository.MutualMatchRepository {
	return &mutualMatchRepository{db: db}
}

func (r *mutualMatchRepository) GetMutualMatch(
	ctx context.Context,
	userID1, userID2 uuid.UUID) (*entity.MutualMatch, error) {

	var mutualMatch entity.MutualMatch

	err := r.db.GormDB.WithContext(ctx).
		Where("user_id_1 = ? AND user_id_2 = ? AND is_deleted = ?",
			userID1, userID2, false).
		First(&mutualMatch).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &mutualMatch, nil
}

func (r *mutualMatchRepository) DeactivateMutualMatchTx(
	ctx context.Context,
	tx *gorm.DB,
	userID1, userID2 uuid.UUID) error {

	now := time.Now().UTC()
	return tx.WithContext(ctx).
		Model(&entity.MutualMatch{}).
		Where("user_id_1 = ? AND user_id_2 = ? AND is_deleted = ?",
			userID1, userID2, false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

// upsert method : create if not exists, reactivate if exists and deleted
func (r *mutualMatchRepository) UpsertMutualMatchTx(
	ctx context.Context,
	tx *gorm.DB,
	userID1, userID2 uuid.UUID) error {

	now := time.Now().UTC()

	mutualMatch := &entity.MutualMatch{
		UserID1:   userID1,
		UserID2:   userID2,
		MatchedAt: now,
		IsDeleted: false,
	}

	return tx.WithContext(ctx).
		Where("user_id_1 = ? AND user_id_2 = ?", userID1, userID2).
		Assign(map[string]interface{}{
			"is_deleted": false,
			"deleted_at": nil,
			"matched_at": now,
			"updated_at": now,
		}).
		FirstOrCreate(mutualMatch).Error
}

func (r *mutualMatchRepository) GetMutualMatchedUserIDs(
	ctx context.Context,
	userID uuid.UUID) ([]uuid.UUID, error) {

	var userIDs []uuid.UUID

	const sql = `
    SELECT
      CASE WHEN user_id_1 = ? THEN user_id_2 ELSE user_id_1 END AS user_id
    FROM mutual_matches
    WHERE (user_id_1 = ? OR user_id_2 = ?) AND is_deleted = false
    ORDER BY updated_at DESC
    `
	// Bind userID three times: once for the CASE, twice for the WHERE clauses
	err := r.db.GormDB.WithContext(ctx).
		Raw(sql, userID, userID, userID).
		Scan(&userIDs).Error

	if err != nil {
		return nil, err
	}

	return userIDs, nil
}
