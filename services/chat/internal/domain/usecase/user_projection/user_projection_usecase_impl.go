package userprojection

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"
)

type userProjectionUsecase struct {
	userProjectionRepository repository.UserProjectionRepository
}

func NewUserProjectionUsecase(userProjectionRepository repository.UserProjectionRepository) usecase.UserProjectionUsecase {
	return &userProjectionUsecase{
		userProjectionRepository: userProjectionRepository,
	}
}

func (u *userProjectionUsecase) CreateUserProjection(
	ctx context.Context,
	userProjection *entity.UserProjection) error {
	return u.userProjectionRepository.CreateUserProjection(ctx, userProjection)
}
