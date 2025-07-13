package auth

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary User login
// @Description User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_login_request body dto.UserLoginRequest true "User login request"
// @Success 200 {object} dto.UserLoginResponse "User login response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/auth/user/login [post]
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := h.authUsecase.UserLogin(c.Request.Context(), req)
	if err != nil {
		log.Printf("User login failed: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User logged in successfully", resp)
}
