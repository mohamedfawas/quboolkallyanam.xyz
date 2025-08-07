package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary User verification
// @Description User verification
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_verification_request body dto.UserVerificationRequest true "User verification request"
// @Success 200 {object} dto.UserVerificationResponse "User verification response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 404 {object} dto.NotFoundError "User not found/OTP expired"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/auth/user/verify [post]
func (h *AuthHandler) UserVerification(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	var req dto.UserVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	user, err := h.authUsecase.UserVerification(reqCtx.Ctx, req, h.config)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to verify user", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User verified successfully")
	apiresponse.Success(c, "User verified successfully", user)
}
