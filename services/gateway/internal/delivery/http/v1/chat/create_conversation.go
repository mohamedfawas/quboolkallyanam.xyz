package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
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
// @Failure 404 {object} apiresponse.Response "User not found"
// @Security BearerAuth
// @Router /api/v1/chat/conversation [post]
func (h *ChatHandler) CreateConversation(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	response, err := h.chatUsecase.CreateConversation(authCtx.Ctx, req)
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
