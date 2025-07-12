package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/user/delete [post]
func (h *AuthHandler) UserDelete(c *gin.Context) {
	var req dto.UserDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	err := h.authUsecase.UserDelete(c.Request.Context(), req)
	if err != nil {
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User deleted successfully", nil)
}
