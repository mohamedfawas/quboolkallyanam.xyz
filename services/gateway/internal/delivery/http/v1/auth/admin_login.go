package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Admin login
// @Description Admin login
// @Tags Auth
// @Accept json
// @Produce json
// @Param admin_login_request body dto.AdminLoginRequest true "Admin login request"
// @Success 200 {object} dto.AdminLoginResponse "Admin login response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized - invalid credentials"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.authUsecase.AdminLogin(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to login", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Admin login successful")
	apiresponse.Success(c, "Admin logged in successfully", resp)
}
