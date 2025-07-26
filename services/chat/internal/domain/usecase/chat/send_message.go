package chat

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *chatUsecase) SendMessage(ctx context.Context, conversationIDStr, senderID, content string) (*entity.Message, error) {
	conversationObjID, err := primitive.ObjectIDFromHex(conversationIDStr)
	if err != nil {
		log.Printf("failed to parse conversation ID: %v", err)
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	now := time.Now().UTC()
	message := &entity.Message{
		ID:             entity.NewMessageID(),
		ConversationID: conversationObjID,
		SenderID:       entity.UserID(senderID),
		Content:        content,
		SentAt:         now,
	}

	if err := c.messageRepository.CreateMessage(ctx, message); err != nil {
		log.Printf("failed to create message: %v", err)
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	log.Printf("Message created successfully: %s", message.ID.Hex())
	return message, nil
}