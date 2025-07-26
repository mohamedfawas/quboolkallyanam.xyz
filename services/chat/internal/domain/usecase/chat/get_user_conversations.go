package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

func (c *chatUsecase) GetUserConversations(ctx context.Context,
	userID uuid.UUID,
	limit,
	offset int) ([]*entity.Conversation, *pagination.PaginationData, error) {

	if limit < 0 {
		limit = constants.DefaultConversationDisplayLimit
	}

	if limit > constants.MaxConversationDisplayLimit {
		limit = constants.MaxConversationDisplayLimit
	}

	if offset < 0 {
		offset = 0
	}

	userIDStr := userID.String()

	conversations, totalCount, err := c.conversationRepository.GetUserConversations(ctx, userIDStr, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	hasMore := (offset + limit) < int(totalCount)

	paginationData := &pagination.PaginationData{
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
	}

	return conversations, paginationData, nil
}
