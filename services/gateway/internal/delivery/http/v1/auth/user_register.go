package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authUsecase.UserRegister(c.Request.Context(), req)
	if err != nil {
		switch err {
		case errors.ErrInvalidEmail:
			apiresponse.Error(c, http.StatusBadRequest, "Email should be in a valid format", errors.ErrInvalidEmail.Error())
			return
		case errors.ErrInvalidPhoneNumber:
			apiresponse.Error(c, http.StatusBadRequest, "Phone number should be in a valid format", errors.ErrInvalidPhoneNumber.Error())
			return
		case errors.ErrInvalidPassword:
			apiresponse.Error(c, http.StatusBadRequest, "Password should be in a valid format", errors.ErrInvalidPassword.Error())
			return
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
			return
		}
	}

	apiresponse.Success(c, http.StatusOK, "OTP sent to the registered email", user)
}
