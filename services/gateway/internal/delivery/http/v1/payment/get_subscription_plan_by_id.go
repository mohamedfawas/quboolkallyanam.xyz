package payment

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

// @Summary Get subscription plan by ID
// @Description Get a specific subscription plan by its ID
// @Tags Payment
// @Accept json
// @Produce json
// @Param plan_id query string true "Subscription plan ID"
// @Success 200 {object} dto.SubscriptionPlan "Subscription plan details"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 404 {object} apiresponse.Response "Plan not found"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/payment/subscription-plan [get]
func (h *PaymentHandler) GetSubscriptionPlanByID(c *gin.Context) {
	planID := c.Query("plan_id")
	if planID == "" {
		log.Printf("Plan ID is required")
		apiresponse.Fail(c, appError.ErrMissingRequiredFields)
		return
	}

	ctx := context.Background()

	plan, err := h.paymentUsecase.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		log.Printf("Failed to get subscription plan with ID %s: %v", planID, err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Subscription plan retrieved successfully", plan)
}
