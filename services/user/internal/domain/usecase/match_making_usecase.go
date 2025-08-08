package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type MatchMakingUsecase interface {
	RecommendUserProfiles(ctx context.Context,
		userID uuid.UUID,
		limit, offset int) ([]*entity.UserProfileResponse, *pagination.PaginationData, error)
	RecordMatchAction(ctx context.Context,
		userID uuid.UUID,
		targetProfileID int64,
		action string) (bool, error)
	GetProfilesByMatchAction(ctx context.Context,
		userID uuid.UUID,
		action string,
		limit, offset int) ([]*entity.UserProfileResponse, *pagination.PaginationData, error)
}
