package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

func (h *PaymentHandler) GetActiveSubscriptionByUserID(c *gin.Context) {
	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		log.Printf("User ID not found in context")
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	subscription, err := h.paymentUsecase.GetActiveSubscriptionByUserID(ctx)
	if err != nil {
		log.Printf("Failed to get active subscription for user %s: %v", userID, err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Active subscription retrieved successfully", subscription)
}
