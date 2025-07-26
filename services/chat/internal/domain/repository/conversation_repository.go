package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type ConversationRepository interface {
	CreateConversation(ctx context.Context,
		conversation *entity.Conversation) error
	GetConversationByParticipants(ctx context.Context,
		participants []string) (*entity.Conversation, error)
	GetConversationByID(ctx context.Context,
		conversationID string) (*entity.Conversation, error)
	GetUserConversations(ctx context.Context,
		userID string,
		limit, offset int) ([]*entity.Conversation, int64, error)
}
