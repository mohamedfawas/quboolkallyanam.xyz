package userprofile

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/event"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type userProfileUsecase struct {
	userProfileRepository repository.UserProfileRepository
	eventPublisher event.EventPublisher
}

func NewUserProfileUsecase(userProfileRepository repository.UserProfileRepository, eventPublisher event.EventPublisher) usecase.UserProfileUsecase {
	return &userProfileUsecase{userProfileRepository: userProfileRepository, eventPublisher: eventPublisher}
}
