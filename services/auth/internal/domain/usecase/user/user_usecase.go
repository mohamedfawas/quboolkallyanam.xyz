package user

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
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
	messageBroker   messaging.Client
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	jwtManager jwt.JWTManager,
	tokenRepository repository.TokenRepository,
	config *config.Config,
	messageBroker messaging.Client,
) usecase.UserUsecase {
	return &userUseCase{
		userRepository:  userRepository,
		jwtManager:      jwtManager,
		tokenRepository: tokenRepository,
		config:          config,
		messageBroker:   messageBroker,
	}
}
