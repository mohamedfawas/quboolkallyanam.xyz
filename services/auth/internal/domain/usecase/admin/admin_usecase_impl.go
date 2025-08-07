package admin

import (
	jwtManager "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type adminUsecase struct {
	adminRepository repository.AdminRepository
	tokenRepository repository.TokenRepository
	userRepository  repository.UserRepository
	jwtManager      jwtManager.JWTManager
	config          *config.Config
}

func NewAdminUsecase(
	adminRepository repository.AdminRepository,
	tokenRepository repository.TokenRepository,
	userRepository repository.UserRepository,
	jwtManager jwtManager.JWTManager,
	config *config.Config) usecase.AdminUsecase {

	return &adminUsecase{
		adminRepository: adminRepository,
		tokenRepository: tokenRepository,
		userRepository:  userRepository,
		jwtManager:      jwtManager,
		config:          config,
	}
}
