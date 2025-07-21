package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Create conversation
// @Description Create a new conversation between users
// @Tags Chat
// @Accept json
// @Produce json
// @Param create_conversation_request body dto.CreateConversationRequest true "Create conversation request"
// @Success 200 {object} dto.CreateConversationResponse "Conversation created successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/chat/conversation [post]
func (h *ChatHandler) CreateConversation(c *gin.Context) {
	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	response, err := h.chatUsecase.CreateConversation(ctx, req)
	if err != nil {
		log.Printf("Failed to create conversation: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Conversation created successfully", response)
}
