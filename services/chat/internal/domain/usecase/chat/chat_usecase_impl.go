package chat

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"
)

type chatUsecase struct {
	conversationRepository   repository.ConversationRepository
	messageRepository        repository.MessageRepository
	userProjectionRepository repository.UserProjectionRepository
}

func NewChatUsecase(
	conversationRepository repository.ConversationRepository,
	messageRepository repository.MessageRepository,
	userProjectionRepository repository.UserProjectionRepository) usecase.ChatUsecase {
	return &chatUsecase{
		conversationRepository:   conversationRepository,
		messageRepository:        messageRepository,
		userProjectionRepository: userProjectionRepository,
	}
}
