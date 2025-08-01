package chat

import (
	"context"
	"fmt"
	"time"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *chatUsecase) SendMessage(ctx context.Context, conversationIDStr, senderID, content string) (*entity.SendMessageResponse, error) {
	conversationObjID, err := primitive.ObjectIDFromHex(conversationIDStr)
	if err != nil {
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
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	userProjection, err := c.userProjectionRepository.GetUserProjectionByUUID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user projection: %w", err)
	}

	if userProjection == nil {
		return nil, appError.ErrUserNotFound
	}

	sendMessageResponse := &entity.SendMessageResponse{
		MessageID:      message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		SenderName:     userProjection.FullName,
		Content:        message.Content,
		SentAt:         message.SentAt,
	}

	return sendMessageResponse, nil
}
