package chat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
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
	requestID, exists := c.Get(constants.ContextKeyRequestID)
	if !exists {
		h.logger.Error("failed to get request ID from context")
		apiresponse.Error(c, fmt.Errorf("request ID context missing"), nil)
		return
	}

	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		h.logger.Error("failed to get user ID from context")
		apiresponse.Error(c, apperrors.ErrUserContextMissing, nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, constants.ContextKeyRequestID, requestID)

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID.(string)),
		zap.String(constants.ContextKeyUserID, userID.(string)),
	)

	response, err := h.chatUsecase.CreateConversation(ctx, req)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("failed to create conversation", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("conversation created successfully")
	apiresponse.Success(c, "Conversation created successfully", response)
}
