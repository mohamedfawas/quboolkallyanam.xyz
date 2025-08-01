package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type ChatUsecase interface {
	CreateConversation(ctx context.Context, req dto.CreateConversationRequest) (*dto.CreateConversationResponse, error)
	SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetConversation(ctx context.Context, req dto.GetConversationRequest) (*dto.GetConversationResponse, error)
	GetMessagesByConversationId(ctx context.Context, req dto.GetMessagesByConversationIdRequest) (*dto.PublicGetMessagesByConversationIdResponse, error)
}
