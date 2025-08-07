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

// @Summary User register
// @Description User register
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_register_request body dto.UserRegisterRequest true "User register request"
// @Success 200 {object} dto.UserRegisterResponse "User register response"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 409 {object} dto.ConflictError "Conflict - email/phone already exists"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/auth/user/register [post]
func (h *AuthHandler) UserRegister(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	user, err := h.authUsecase.UserRegister(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to register user", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User register request processed successfully")
	apiresponse.Success(c, "OTP sent to the registered email", user)
}
