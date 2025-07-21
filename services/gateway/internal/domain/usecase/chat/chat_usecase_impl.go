package chat

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type chatUsecase struct {
	chatClient client.ChatClient
}

func NewChatUsecase(chatClient client.ChatClient) usecase.ChatUsecase {
	return &chatUsecase{chatClient: chatClient}
}
