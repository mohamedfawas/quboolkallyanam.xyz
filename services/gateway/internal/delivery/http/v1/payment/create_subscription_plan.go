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

func (h *PaymentHandler) CreateSubscriptionPlan(c *gin.Context) {
	var req dto.CreateSubscriptionPlanRequest
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

	updateReq := dto.UpdateSubscriptionPlanRequest{
		ID:           req.ID,
		DurationDays: &req.DurationDays,
		Amount:       &req.Amount,
		Currency:     &req.Currency,
		Description:  &req.Description,
		IsActive:     &req.IsActive,
	}

	response, err := h.paymentUsecase.CreateOrUpdateSubscriptionPlan(ctx, updateReq)
	if err != nil {
		log.Printf("Failed to create subscription plan: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Subscription plan created successfully", response)
}