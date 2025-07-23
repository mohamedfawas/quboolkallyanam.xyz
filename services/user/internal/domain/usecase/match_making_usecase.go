package usecase

import (
	"context"

	"github.com/google/uuid"
)

type MatchMakingUsecase interface {
	// GetMatches(ctx context.Context, userID uuid.UUID) ([]entity.UserProfile, error)
	RecordMatchAction(ctx context.Context, userID uuid.UUID, targetProfileID uint, action string) (bool, error)
}
