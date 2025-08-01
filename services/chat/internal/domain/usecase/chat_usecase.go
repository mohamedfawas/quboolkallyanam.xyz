package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type ChatUsecase interface {
	CreateConversation(ctx context.Context,
		userUUIDStr string,
		partnerProfileID int64) (*entity.CreateConversationResponse, error)
	SendMessage(ctx context.Context,
		conversationID,
		senderID,
		content string) (*entity.SendMessageResponse, error)
	GetConversationByID(ctx context.Context,
		conversationID string) (*entity.GetConversationResponse, error)
	GetMessagesByConversationID(ctx context.Context,
		conversationID string,
		userID string,
		limit, offset int32) (*entity.GetMessagesByConversationIDResponse, error)
}
