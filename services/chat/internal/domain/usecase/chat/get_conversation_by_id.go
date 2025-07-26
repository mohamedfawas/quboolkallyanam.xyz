package chat

import (
	"context"
	"log"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

func (c *chatUsecase) GetConversationByID(ctx context.Context, 
	conversationID string) (*entity.Conversation, error) {

	conversation, err := c.conversationRepository.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get conversation by ID: %v", err)
		return nil, err
	}

	return conversation, nil
}