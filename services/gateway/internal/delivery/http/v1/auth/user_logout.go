package auth

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

// @Summary User logout
// @Description User logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_logout_request body dto.UserLogoutRequest true "User logout request"
// @Success 200 {object} dto.UserLogoutResponse "User logout response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/user/logout [post]
func (h *AuthHandler) UserLogout(c *gin.Context) {
	accessToken, exists := c.Get(constants.ContextKeyAccessToken)
	if !exists {
		log.Printf("Access token not found in context")
		apiresponse.Fail(c, fmt.Errorf("access token not found in context"))
		return
	}

	err := h.authUsecase.UserLogout(c.Request.Context(), accessToken.(string))
	if err != nil {
		log.Printf("User logout failed: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User logged out successfully", nil)
}
