package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type ChatUsecase interface {
	CreateConversation(ctx context.Context, userUUIDStr string, partnerProfileID int64) (*entity.Conversation, error)
	SendMessage(ctx context.Context, conversationID, senderID, content string) (*entity.Message, error)
}
