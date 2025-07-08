package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
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
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/admin/logout [post]
func (h *AuthHandler) AdminLogout(c *gin.Context) {
	var req dto.AdminLogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := h.authUsecase.AdminLogout(c.Request.Context(), req)
	if err != nil {
		apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
		return
	}

	apiresponse.Success(c, http.StatusOK, "Admin logged out successfully", nil)
}
