package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	resp, err := h.authUsecase.RefreshToken(c.Request.Context(), req)
	if err != nil {
		apiresponse.Error(c, http.StatusInternalServerError, "Something went wrong. Please try again later.", errors.ErrInternalServerError.Error())
		return
	}

	apiresponse.Success(c, http.StatusOK, "Token refreshed successfully", resp)
}
