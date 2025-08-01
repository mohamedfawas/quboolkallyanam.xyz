package chat

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *chatUsecase) GetConversation(ctx context.Context, req dto.GetConversationRequest) (*dto.GetConversationResponse, error) {
	response, err := u.chatClient.GetConversation(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
