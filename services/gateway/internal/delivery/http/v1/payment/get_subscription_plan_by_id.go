package payment

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

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
