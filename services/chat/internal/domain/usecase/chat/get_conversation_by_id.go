package chat

import (
	"context"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

func (c *chatUsecase) GetConversationByID(ctx context.Context,
	conversationID string) (*entity.GetConversationResponse, error) {

	conversation, err := c.conversationRepository.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, appError.ErrConversationNotFound
	}

	participantIDs := make([]string, len(conversation.ParticipantIDs))
	for i, id := range conversation.ParticipantIDs {
		participantIDs[i] = string(id)
	}

	getConversationResponse := &entity.GetConversationResponse{
		ConversationID: conversation.ID,
		ParticipantIDs: participantIDs,
		CreatedAt:      conversation.CreatedAt,
		UpdatedAt:      conversation.UpdatedAt,
	}

	return getConversationResponse, nil
}
