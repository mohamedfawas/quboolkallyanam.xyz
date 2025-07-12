package user

import (
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
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
	messageBroker   messageBroker.Client
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	jwtManager jwt.JWTManager,
	tokenRepository repository.TokenRepository,
	config *config.Config,
	messageBroker messageBroker.Client,
) usecase.UserUsecase {
	return &userUseCase{
		userRepository:  userRepository,
		jwtManager:      jwtManager,
		tokenRepository: tokenRepository,
		config:          config,
		messageBroker:   messageBroker,
	}
}
