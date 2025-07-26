package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"gorm.io/gorm"
)

type MutualMatchRepository interface {
	GetMutualMatch(ctx context.Context, userID1, userID2 uuid.UUID) (*entity.MutualMatch, error)
	DeactivateMutualMatchTx(ctx context.Context, tx *gorm.DB, userID1, userID2 uuid.UUID) error
	UpsertMutualMatchTx(ctx context.Context, tx *gorm.DB, userID1, userID2 uuid.UUID) error
	GetMutualMatchedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
