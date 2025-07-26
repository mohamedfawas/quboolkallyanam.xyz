package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"

	"github.com/google/uuid"
)

type MatchMakingUsecase interface {
	RecommendUserProfiles(ctx context.Context,
		userID uuid.UUID,
		limit, offset int) ([]*entity.UserProfileResponse, *entity.PaginationData, error)
	RecordMatchAction(ctx context.Context,
		userID uuid.UUID,
		targetProfileID uint, action string) (bool, error)
	GetProfilesByMatchAction(ctx context.Context,
		userID uuid.UUID,
		action string,
		limit, offset int) ([]*entity.UserProfileResponse, *entity.PaginationData, error)
}
