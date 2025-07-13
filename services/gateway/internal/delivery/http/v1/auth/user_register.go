package auth

import (
	"log"

	"github.com/gin-gonic/gin"
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
	ctx := c.Request.Context()

	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	user, err := h.authUsecase.UserRegister(ctx, req)
	if err != nil {
		log.Printf("User registration failed: %v", err)
		apiresponse.Fail(c, err)
		return
	}
	apiresponse.Success(c, "OTP sent to the registered email", user)
}
