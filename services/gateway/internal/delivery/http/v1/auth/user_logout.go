package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary User logout
// @Description User logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserLogoutResponse "User logout response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/user/logout [post]
func (h *AuthHandler) UserLogout(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	accessToken, exists := c.Get(constants.ContextKeyAccessToken)
	if !exists {
		apiresponse.Error(c, apperrors.ErrAccessTokenNotFound, nil)
		return
	}

	err = h.authUsecase.UserLogout(authCtx.Ctx, accessToken.(string))
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to logout", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User logged out successfully")
	apiresponse.Success(c, "User logged out successfully", nil)
}
