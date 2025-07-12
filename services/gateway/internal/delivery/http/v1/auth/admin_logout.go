package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Admin logout
// @Description Admin logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param admin_logout_request body dto.AdminLogoutRequest true "Admin logout request"
// @Success 200 {object} dto.AdminLogoutResponse "Admin logout response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/admin/logout [post]
func (h *AuthHandler) AdminLogout(c *gin.Context) {
	var req dto.AdminLogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	err := h.authUsecase.AdminLogout(c.Request.Context(), req)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Admin logged out successfully", nil)
}
