package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
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
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	resp, err := h.authUsecase.AdminLogin(c.Request.Context(), req)
	if err != nil {
		switch err {
		case errors.ErrInvalidEmail:
			apiresponse.Error(c, http.StatusBadRequest, "Email should be in a valid format", errors.ErrInvalidEmail.Error())
			return
		case errors.ErrInvalidPassword:
			apiresponse.Error(c, http.StatusBadRequest, "Password should be in a valid format", errors.ErrInvalidPassword.Error())
			return
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
			return
		}
	}

	apiresponse.Success(c, http.StatusOK, "Admin logged in successfully", resp)
}
