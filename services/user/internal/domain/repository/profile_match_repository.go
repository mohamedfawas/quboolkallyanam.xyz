package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"gorm.io/gorm"
)

type ProfileMatchRepository interface {
	GetExistingMatch(ctx context.Context, userID uuid.UUID, targetID uuid.UUID) (*entity.ProfileMatch, error)
	CreateMatchAction(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error
	CreateMatchActionTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error
	UpdateMatchAction(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error
	UpdateMatchActionTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, targetID uuid.UUID, isLiked bool) error
}
