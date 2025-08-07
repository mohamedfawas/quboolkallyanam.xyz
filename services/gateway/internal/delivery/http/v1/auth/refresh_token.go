package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Refresh token
// @Description Refresh the access token using a valid refresh token passed in headers
// @Tags Auth
// @Accept json
// @Produce json
// @Param Refresh-Token header string true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse "Refresh token response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized - invalid credentials"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/auth/user/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	refreshToken := c.GetHeader(constants.HeaderRefreshToken)
	if strings.TrimSpace(refreshToken) == "" {
		apiresponse.Error(c, apperrors.ErrRefreshTokenNotFound, nil)
		return
	}

	req := dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := h.authUsecase.RefreshToken(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to refresh token", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Token refreshed successfully")
	apiresponse.Success(c, "Token refreshed successfully", resp)
}
