package chat

import (

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type ChatHandler struct {
	chatUsecase usecase.ChatUsecase
}

func NewChatHandler(chatUsecase usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{chatUsecase: chatUsecase}
}

