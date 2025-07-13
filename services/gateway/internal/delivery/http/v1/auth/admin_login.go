package auth

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Admin login
// @Description Admin login
// @Tags Auth
// @Accept json
// @Produce json
// @Param admin_login_request body dto.AdminLoginRequest true "Admin login request"
// @Success 200 {object} dto.AdminLoginResponse "Admin login response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := h.authUsecase.AdminLogin(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed to login: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Admin logged in successfully", resp)
}
