package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type UserProjectionRepository interface {
	CreateUserProjection(ctx context.Context, userProjection *entity.UserProjection) error
	GetUserProjectionByID(ctx context.Context, id int64) (*entity.UserProjection, error)
}
