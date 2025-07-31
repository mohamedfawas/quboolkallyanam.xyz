package user

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	logger      *zap.Logger
}

func NewUserHandler(userUsecase usecase.UserUsecase,
	logger *zap.Logger) *UserHandler {

	return &UserHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}
