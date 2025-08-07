package payment

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

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
