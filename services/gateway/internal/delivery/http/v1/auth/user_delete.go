package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary User delete
// @Description User delete
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_delete_request body dto.UserDeleteRequest true "User delete request"
// @Success 200 {object} dto.UserDeleteResponse "User delete response"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/user/delete [post]
func (h *AuthHandler) UserDelete(c *gin.Context) {
	var req dto.UserDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := h.authUsecase.UserDelete(c.Request.Context(), req)
	if err != nil {
		switch err {
		case errors.ErrInvalidPassword:
			apiresponse.Error(c, http.StatusBadRequest, "Password should contain at least 8 characters, one capital letter,one number, one small letter and one special character", errors.ErrInvalidPassword.Error())
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
		}
		return
	}

	apiresponse.Success(c, http.StatusOK, "User deleted successfully", nil)
}
