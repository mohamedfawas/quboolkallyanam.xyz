package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
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
	var req dto.UserLogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	err := h.authUsecase.UserLogout(c.Request.Context(), req)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User logged out successfully", nil)
}
