package payment

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

// @Summary Get active subscription plans
// @Description Get all active subscription plans
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {array} dto.SubscriptionPlan "List of active subscription plans"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/payment/subscription-plans [get]
func (h *PaymentHandler) GetActiveSubscriptionPlans(c *gin.Context) {
	ctx := context.Background()

	plans, err := h.paymentUsecase.GetActiveSubscriptionPlans(ctx)
	if err != nil {
		log.Printf("Failed to get active subscription plans: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Active subscription plans retrieved successfully", plans)
}
