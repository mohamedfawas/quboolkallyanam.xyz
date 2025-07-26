package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type ChatUsecase interface {
	CreateConversation(ctx context.Context,
		userUUIDStr string,
		partnerProfileID int64) (*entity.Conversation, error)
	SendMessage(ctx context.Context,
		conversationID,
		senderID,
		content string) (*entity.Message, error)
	GetConversationByID(ctx context.Context,
		conversationID string) (*entity.Conversation, error)
	GetUserConversations(ctx context.Context,
		userID uuid.UUID,
		limit,
		offset int) ([]*entity.Conversation, *pagination.PaginationData, error)
}
