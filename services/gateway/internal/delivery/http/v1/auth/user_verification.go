package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *AuthHandler) UserVerification(c *gin.Context) {
	var req dto.UserVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authUsecase.UserVerification(c.Request.Context(), req, h.config)
	if err != nil {
		switch err {
		case errors.ErrInvalidEmail:
			apiresponse.Error(c, http.StatusBadRequest, "Email should be in a valid format", errors.ErrInvalidEmail.Error())
			return
		case errors.ErrInvalidOTP:
			apiresponse.Error(c, http.StatusBadRequest, "OTP should be in a valid format", errors.ErrInvalidOTP.Error())
			return
		default:
			apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
			return
		}
	}

	apiresponse.Success(c, http.StatusOK, "User verified successfully", user)
}
