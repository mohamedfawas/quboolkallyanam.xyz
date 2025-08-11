package userprofile

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/event"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type userProfileUsecase struct {
	userProfileRepository        repository.UserProfileRepository
	userImageRepository          repository.UserImageRepository
	partnerPreferencesRepository repository.PartnerPreferencesRepository
	eventPublisher               event.EventPublisher
	photoStorage                 mediastorage.PhotoStorage
	config                       *config.Config
}

func NewUserProfileUsecase(
	userProfileRepository repository.UserProfileRepository,
	userImageRepository repository.UserImageRepository,
	partnerPreferencesRepository repository.PartnerPreferencesRepository,
	eventPublisher event.EventPublisher,
	photoStorage mediastorage.PhotoStorage,
	config *config.Config,
) usecase.UserProfileUsecase {
	return &userProfileUsecase{
		userProfileRepository:        userProfileRepository,
		userImageRepository:          userImageRepository,
		partnerPreferencesRepository: partnerPreferencesRepository,
		eventPublisher:               eventPublisher,
		photoStorage:                 photoStorage,
		config:                       config,
	}
}
