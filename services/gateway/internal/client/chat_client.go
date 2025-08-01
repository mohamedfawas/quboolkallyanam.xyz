package client

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type ChatClient interface {
	CreateConversation(ctx context.Context, req dto.CreateConversationRequest) (*dto.CreateConversationResponse, error)
	SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetConversation(ctx context.Context, req dto.GetConversationRequest) (*dto.GetConversationResponse, error)
	GetMessagesByConversationId(ctx context.Context, req dto.GetMessagesByConversationIdRequest) (*dto.GetMessagesByConversationIdResponse, error)
}
