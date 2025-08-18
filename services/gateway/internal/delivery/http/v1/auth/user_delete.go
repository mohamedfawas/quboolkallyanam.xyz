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

// @Summary User delete
// @Description User delete
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_delete_request body dto.UserDeleteRequest true "User delete request"
// @Success 200 {object} dto.UserDeleteResponse "User delete response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized - invalid credentials"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/user/delete [post]
func (h *AuthHandler) UserDelete(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.UserDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	err = h.authUsecase.UserDelete(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to delete user", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	metrics.UsersDeletedTotal.Inc()
	log.Info("User deleted successfully")
	apiresponse.Success(c, "User deleted successfully", nil)
}
