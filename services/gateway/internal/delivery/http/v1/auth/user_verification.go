package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary User verification
// @Description User verification
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_verification_request body dto.UserVerificationRequest true "User verification request"
// @Success 200 {object} dto.UserVerificationResponse "User verification response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/auth/user/verify [post]
func (h *AuthHandler) UserVerification(c *gin.Context) {
	var req dto.UserVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	user, err := h.authUsecase.UserVerification(c.Request.Context(), req, h.config)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User verified successfully", user)
}
