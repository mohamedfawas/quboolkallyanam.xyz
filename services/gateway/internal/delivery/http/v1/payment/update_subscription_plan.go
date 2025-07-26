package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *PaymentHandler) UpdateSubscriptionPlan(c *gin.Context) {
	var req dto.UpdateSubscriptionPlanRequest
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

	response, err := h.paymentUsecase.CreateOrUpdateSubscriptionPlan(ctx, req)
	if err != nil {
		log.Printf("Failed to update subscription plan: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Subscription plan updated successfully", response)
}