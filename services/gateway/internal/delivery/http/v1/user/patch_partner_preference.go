package user

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *UserHandler) PatchPartnerPreference(c *gin.Context) {
	var req dto.PartnerPreferencePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	err := h.userUsecase.UpdateUserPartnerPreferences(ctx, constants.UpdateOperationType, req)
	if err != nil {
		log.Printf("Failed to update partner preference: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Partner preference updated successfully", nil)
}
