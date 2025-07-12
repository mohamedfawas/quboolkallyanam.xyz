package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Refresh token
// @Description Refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token_request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.RefreshTokenResponse "Refresh token response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := h.authUsecase.RefreshToken(c.Request.Context(), req)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Token refreshed successfully", resp)
}
