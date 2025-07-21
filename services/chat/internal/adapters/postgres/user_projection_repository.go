package postgres

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
)

type userProjectionRepository struct {
	db *postgres.Client
}

func NewUserProjectionRepository(db *postgres.Client) repository.UserProjectionRepository {
	return &userProjectionRepository{db: db}
}

func (r *userProjectionRepository) CreateUserProjection(ctx context.Context, userProjection *entity.UserProjection) error {
	return r.db.GormDB.WithContext(ctx).Create(userProjection).Error
}

func (r *userProjectionRepository) GetUserProjectionByID(ctx context.Context, id int64) (*entity.UserProjection, error) {
	var userProjection entity.UserProjection
	if err := r.db.GormDB.WithContext(ctx).Where("user_profile_id = ?", id).First(&userProjection).Error; err != nil {
		return nil, err
	}
	return &userProjection, nil
}
