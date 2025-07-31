package userprojection

import (
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
