package admin

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/event"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type adminUsecase struct {
	adminRepository repository.AdminRepository
	tokenRepository repository.TokenRepository
	userRepository  repository.UserRepository
	jwtManager      jwt.JWTManager
	eventPublisher  event.EventPublisher
	config          *config.Config
}

func NewAdminUsecase(
	adminRepository repository.AdminRepository,
	tokenRepository repository.TokenRepository,
	userRepository repository.UserRepository,
	jwtManager jwt.JWTManager,
	eventPublisher event.EventPublisher,
	config *config.Config) usecase.AdminUsecase {

	return &adminUsecase{
		adminRepository: adminRepository,
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
		jwtManager:      jwtManager,
		eventPublisher:  eventPublisher,
		config:          config,
	}
}
