package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type UserProjectionUsecase interface {
	CreateOrUpdateUserProjection(ctx context.Context, userProjection *entity.UserProjection) error
}
