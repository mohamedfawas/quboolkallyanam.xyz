package chat

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Get user conversations
// @Description Get paginated list of conversations for the authenticated user
// @Tags Chat
// @Accept json
// @Produce json
// @Param limit query int false "Number of conversations to return (default: 10, max: 50)"
// @Param offset query int false "Number of conversations to skip (default: 0)"
// @Success 200 {object} dto.GetUserConversationsResponse "User conversations retrieved successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/chat/conversations [get]
func (h *ChatHandler) GetUserConversations(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		log.Printf("Invalid limit parameter: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid limit parameter: %w", err))
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		log.Printf("Invalid offset parameter: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid offset parameter: %w", err))
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	req := dto.GetUserConversationsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	response, err := h.chatUsecase.GetUserConversations(ctx, req)
	if err != nil {
		log.Printf("Failed to get user conversations: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User conversations retrieved successfully", response)
}