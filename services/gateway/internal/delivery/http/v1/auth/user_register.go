package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	logger "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary User register
// @Description User register
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_register_request body dto.UserRegisterRequest true "User register request"
// @Success 200 {object} dto.UserRegisterResponse "User register response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 409 {object} apiresponse.Response "Conflict"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/auth/user/register [post]
func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}
	logger.Log.Info("🔑 UserRegister request validated in handler : ", "email : ", req.Email, "phone : ", req.Phone)

	logger.Log.Info("🔑 UserRegister request sent to usecase from handler : ", "email : ", req.Email, "phone : ", req.Phone)
	user, err := h.authUsecase.UserRegister(c.Request.Context(), req)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	logger.Log.Info("🔑 UserRegister response sent from usecase to handler : ", "email : ", user.Email, "phone : ", user.Phone)
	apiresponse.Success(c, "OTP sent to the registered email", user)
}
