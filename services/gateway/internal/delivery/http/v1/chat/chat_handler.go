package chat

import (
	"go.uber.org/zap"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type ChatHandler struct {
	chatUsecase usecase.ChatUsecase
	logger      *zap.Logger
	jwtManager  *jwt.JWTManager
	connManager *ConnectionManager
}

func NewChatHandler(chatUsecase usecase.ChatUsecase, logger *zap.Logger, jwtManager *jwt.JWTManager) *ChatHandler {
	return &ChatHandler{
		chatUsecase: chatUsecase,
		logger:      logger,
		jwtManager:  jwtManager,
		connManager: NewConnectionManager(logger),
	}
}

