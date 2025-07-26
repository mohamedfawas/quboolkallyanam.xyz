package chat

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *chatUsecase) SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	response, err := u.chatClient.SendMessage(ctx, req)
	if err != nil {
		log.Printf("Failed to send message via chat client: %v", err)
		return nil, err
	}

	return response, nil
}