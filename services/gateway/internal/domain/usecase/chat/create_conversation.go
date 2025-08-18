package chat

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *chatUsecase) CreateConversation(ctx context.Context, req dto.CreateConversationRequest) (*dto.CreateConversationResponse, error) {
	if req.PartnerUserProfileID < 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return u.chatClient.CreateConversation(ctx, req)
}
