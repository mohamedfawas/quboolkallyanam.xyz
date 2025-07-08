package auth

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	config      config.Config
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, config config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		config:      config,
	}
}
