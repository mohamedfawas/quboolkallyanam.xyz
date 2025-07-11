package user

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type userUseCase struct {
	userRepository  repository.UserRepository
	jwtManager      jwt.JWTManager
	tokenRepository repository.TokenRepository
	config          *config.Config
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	jwtManager jwt.JWTManager,
	tokenRepository repository.TokenRepository,
	config *config.Config,
) usecase.UserUsecase {
	return &userUseCase{
		userRepository:  userRepository,
		jwtManager:      jwtManager,
		tokenRepository: tokenRepository,
		config:          config,
	}
}
