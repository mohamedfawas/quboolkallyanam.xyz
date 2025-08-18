package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/metrics"
	"go.uber.org/zap"
)

// @Summary Unblock user
// @Description Unblock a user by email, phone, or ID
// @Tags Auth
// @Accept json
// @Produce json
// @Param unblock_user_request body dto.BlockOrUnblockUserRequest true "Unblock user request"
// @Success 200 {object} dto.BlockOrUnblockUserResponse "Unblock user response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized - invalid credentials"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/admin/unblock-user [post]
func (h *AuthHandler) AdminUnBlockUser(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	var req dto.BlockOrUnblockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	req.ShouldBlock = false

	resp, err := h.authUsecase.BlockOrUnblockUser(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to unblock user", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	metrics.AdminUserUnblockedTotal.Inc()
	log.Info("User unblock request processed successfully",
		zap.String("field", req.Field),
		zap.String("value", req.Value),
	)
	apiresponse.Success(c, "User unblock request processed successfully", resp)
}
