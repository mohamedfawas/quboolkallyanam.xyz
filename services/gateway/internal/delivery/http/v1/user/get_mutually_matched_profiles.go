package user

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *UserHandler) GetMutuallyMatchedProfiles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		log.Printf("Invalid limit parameter: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid limit parameter: %w", err))
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		log.Printf("Invalid offset parameter: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid offset parameter: %w", err))
		return
	}

	if limit <= 0 || limit > 50 {
		apiresponse.Fail(c, fmt.Errorf("limit must be between 1 and 50"))
		return
	}

	if offset < 0 {
		apiresponse.Fail(c, fmt.Errorf("offset must be non-negative"))
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	request := dto.GetProfilesByMatchActionRequest{
		Action: validation.MatchMakingOptionMutual,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	response, err := h.userUsecase.GetProfilesByMatchAction(ctx, request)
	if err != nil {
		log.Printf("Failed to get mutually matched profiles: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Mutually matched profiles retrieved successfully", response)
}