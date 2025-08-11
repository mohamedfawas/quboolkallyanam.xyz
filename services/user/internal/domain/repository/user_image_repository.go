package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserImageRepository interface {
	DisplayOrderOccupied(ctx context.Context, userID uuid.UUID, displayOrder int32) (bool, error)
	CreateUserImage(ctx context.Context, userImage *entity.UserImage) error
	DeleteUserImage(ctx context.Context, userID uuid.UUID, displayOrder int32) error
	GetUserImage(ctx context.Context, userID uuid.UUID, displayOrder int32) (*entity.UserImage, error)
	ListUserImages(ctx context.Context, userID uuid.UUID) ([]entity.UserImage, error)
}
