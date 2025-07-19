package userprofile

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type userProfileUsecase struct {
	userProfileRepository repository.UserProfileRepository
}

func NewUserProfileUsecase(userProfileRepository repository.UserProfileRepository) usecase.UserProfileUsecase {
	return &userProfileUsecase{userProfileRepository: userProfileRepository}
}
