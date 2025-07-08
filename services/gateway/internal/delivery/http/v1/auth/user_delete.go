package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

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
