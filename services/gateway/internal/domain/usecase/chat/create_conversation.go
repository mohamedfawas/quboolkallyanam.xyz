package chat

import (
	"context"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *chatUsecase) CreateConversation(ctx context.Context, req dto.CreateConversationRequest) (*dto.CreateConversationResponse, error) {
	if req.PartnerUserProfileID < 0 {
		return nil, appErrors.ErrInvalidInput
	}

	return u.chatClient.CreateConversation(ctx, req)
}
