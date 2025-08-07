package auth

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	config      config.Config
	logger      *zap.Logger
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, config config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		config:      config,
		logger:      logger,
	}
}
