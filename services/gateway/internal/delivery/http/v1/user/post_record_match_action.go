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

// @Summary Record match action
// @Description Record a match action (like, dislike, etc.)
// @Tags User
// @Accept json
// @Produce json
// @Param record_match_action_request body dto.RecordMatchActionRequest true "Record match action request"
// @Success 200 {object} dto.RecordMatchActionResponse "Match action recorded successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/match-action [post]
func (h *UserHandler) PostRecordMatchAction(c *gin.Context) {
	var req dto.RecordMatchActionRequest
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

	response, err := h.userUsecase.RecordMatchAction(ctx, req)
	if err != nil {
		log.Printf("Failed to record match action: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Match action recorded successfully", response)
}