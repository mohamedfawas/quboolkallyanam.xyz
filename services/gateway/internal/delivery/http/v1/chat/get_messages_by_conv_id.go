package chat

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Get messages by conversation ID
// @Description Get messages by conversation ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param conversation_id path string true "Conversation ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.PublicGetMessagesByConversationIdResponse "Messages fetched successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Failure 404 {object} apiresponse.Response "Conversation not found"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Security BearerAuth
// @Router /api/v1/chat/conversation/{conversation_id}/messages [get]

func (h *ChatHandler) GetMessagesByConversationId(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		apiresponse.Error(c, apperrors.ErrInvalidConversationID, nil)
		return
	}

	convLogger := log.With(
		zap.String("conversation_id", conversationID),
	)

	var limit, offset int32 = 10, 0 // Default values

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && parsedOffset >= 0 {
			offset = int32(parsedOffset)
		}
	}

	req := dto.GetMessagesByConversationIdRequest{
		ConversationID: conversationID,
		Limit:          limit,
		Offset:         offset,
	}

	response, err := h.chatUsecase.GetMessagesByConversationId(authCtx.Ctx, req)
	if err != nil {
		if !apperrors.IsAppError(err) {
			convLogger.Error("Failed to get messages by conversation ID", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	convLogger.Info("Messages fetched successfully")

	apiresponse.Success(c, "Messages fetched successfully", response)
}
