package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Admin logout
// @Description Admin logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.AdminLogoutResponse "Admin logout response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized - invalid credentials"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/admin/logout [post]
func (h *AuthHandler) AdminLogout(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	accessToken, exists := c.Get(constants.ContextKeyAccessToken)
	if !exists {
		apiresponse.Error(c, apperrors.ErrAccessTokenNotFound, nil)
		return
	}

	err = h.authUsecase.AdminLogout(reqCtx.Ctx, accessToken.(string))
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to logout", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Admin logged out successfully")
	apiresponse.Success(c, "Admin logged out successfully", nil)
}
