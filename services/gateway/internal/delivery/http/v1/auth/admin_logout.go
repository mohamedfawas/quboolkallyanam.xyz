package auth

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
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
	accessToken, exists := c.Get(constants.ContextKeyAccessToken)
	if !exists {
		log.Printf("Access token not found in context")
		apiresponse.Fail(c, fmt.Errorf("access token not found in context"))
		return
	}

	err := h.authUsecase.AdminLogout(c.Request.Context(), accessToken.(string))
	if err != nil {
		log.Printf("Failed to logout: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Admin logged out successfully", nil)
}
