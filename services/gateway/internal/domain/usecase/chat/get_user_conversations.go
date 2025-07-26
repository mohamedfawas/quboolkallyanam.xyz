package chat

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *chatUsecase) GetUserConversations(ctx context.Context, req dto.GetUserConversationsRequest) (*dto.GetUserConversationsResponse, error) {
	response, err := u.chatClient.GetUserConversations(ctx, req)
	if err != nil {
		log.Printf("Failed to get user conversations via chat client: %v", err)
		return nil, err
	}

	return response, nil
}