package chat

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (c *chatUsecase) GetMessagesByConversationId(
	ctx context.Context, 
	req dto.GetMessagesByConversationIdRequest) (*dto.PublicGetMessagesByConversationIdResponse, error) {

	if !validation.IsValidConversationID(req.ConversationID) {
		return nil, apperrors.ErrInvalidConversationID
	}

	internalResponse, err := c.chatClient.GetMessagesByConversationId(ctx, req)
	if err != nil {
		return nil, err
	}

	publicMessages := make([]dto.PublicMessageInfo, 0, len(internalResponse.Messages))
	for _, msg := range internalResponse.Messages {
		publicMessage := dto.PublicMessageInfo{
			MessageID:  msg.MessageID,
			SenderName: msg.SenderName,
			Content:    msg.Content,
			SentAt:     msg.SentAt,
		}
		publicMessages = append(publicMessages, publicMessage)
	}

	publicResponse := &dto.PublicGetMessagesByConversationIdResponse{
		Messages:   publicMessages,
		Pagination: internalResponse.Pagination,
	}

	return publicResponse, nil
}