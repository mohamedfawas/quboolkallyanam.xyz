package chat

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

func (c *chatUsecase) GetMessagesByConversationID(ctx context.Context,
	conversationID string,
	userID string,
	limit, offset int32) (*entity.GetMessagesByConversationIDResponse, error) {
	if limit <= 0 {
		limit = constants.DefaultConversationDisplayLimit
	}

	if limit > constants.MaxConversationDisplayLimit {
		limit = constants.MaxConversationDisplayLimit
	}

	if offset < 0 {
		offset = 0
	}

	if !validation.IsValidConversationID(conversationID) {
		return nil, apperrors.ErrInvalidConversationID
	}

	conversation, err := c.conversationRepository.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation == nil {
		return nil, apperrors.ErrConversationNotFound
	}

	isParticipant := false
	for _, participantID := range conversation.ParticipantIDs {
		if string(participantID) == userID {
			isParticipant = true
			break
		}
	}
	
	if !isParticipant {
		return nil, apperrors.ErrUserNotParticipant
	}

	

	// Create participant name mapping for efficient lookup
	participantNameMap := make(map[string]string)
	for _, participantID := range conversation.ParticipantIDs {
		userProjection, err := c.userProjectionRepository.GetUserProjectionByUUID(ctx, string(participantID))
		if err != nil {
			return nil, fmt.Errorf("failed to get user projection for participant %s: %w", participantID, err)
		}
		if userProjection == nil {
			return nil, apperrors.ErrUserNotFound
		}
		if userProjection != nil {
			participantNameMap[string(participantID)] = userProjection.FullName
		}
	}

	messages, totalCount, err := c.messageRepository.GetMessagesByConversationID(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Convert messages to response format with participant names
	messageInfos := make([]entity.MessageInfo, 0, len(messages))
	for _, message := range messages {
		senderName := participantNameMap[string(message.SenderID)]
		if senderName == "" {
			senderName = "Unknown User" // Fallback name
		}

		messageInfo := entity.MessageInfo{
			MessageID:  message.ID.Hex(),
			SenderID:   string(message.SenderID),
			SenderName: senderName,
			Content:    message.Content,
			SentAt:     message.SentAt,
		}
		messageInfos = append(messageInfos, messageInfo)
	}

	hasMore := (int64(offset) + int64(limit)) < totalCount
	pagination := entity.PaginationInfo{
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
	}

	response := &entity.GetMessagesByConversationIDResponse{
		Messages:   messageInfos,
		Pagination: pagination,
	}

	return response, nil
}
