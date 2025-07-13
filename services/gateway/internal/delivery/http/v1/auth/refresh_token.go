package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
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
	refreshToken := c.GetHeader(constants.HeaderRefreshToken)
	if strings.TrimSpace(refreshToken) == "" {
		log.Printf("Refresh token not found in header")
		apiresponse.Fail(c, fmt.Errorf("refresh token not found in header"))
		return
	}

	req := dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := h.authUsecase.RefreshToken(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Token refreshed successfully", resp)
}
