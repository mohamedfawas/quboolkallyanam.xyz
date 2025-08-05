package matchmaking

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type matchMakingUsecase struct {
	userProfileRepository        repository.UserProfileRepository
	partnerPreferencesRepository repository.PartnerPreferencesRepository
	profileMatchRepository       repository.ProfileMatchRepository
	mutualMatchRepository        repository.MutualMatchRepository
	transactionManager           *postgres.TransactionManager
	photoStorage                 mediastorage.PhotoStorage
	config                       *config.Config
}

func NewMatchMakingUsecase(
	userProfileRepository repository.UserProfileRepository,
	partnerPreferencesRepository repository.PartnerPreferencesRepository,
	profileMatchRepository repository.ProfileMatchRepository,
	mutualMatchRepository repository.MutualMatchRepository,
	transactionManager *postgres.TransactionManager,
	photoStorage mediastorage.PhotoStorage,
	config *config.Config,
) usecase.MatchMakingUsecase {
	return &matchMakingUsecase{
		userProfileRepository:        userProfileRepository,
		partnerPreferencesRepository: partnerPreferencesRepository,
		profileMatchRepository:       profileMatchRepository,
		mutualMatchRepository:        mutualMatchRepository,
		transactionManager:           transactionManager,
		photoStorage:                 photoStorage,
		config:                       config,
	}
}
