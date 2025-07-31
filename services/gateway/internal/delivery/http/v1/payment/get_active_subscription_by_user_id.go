package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

// @Summary Get active subscription by user ID
// @Description Get the active subscription for the authenticated user
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ActiveSubscription "Active subscription details"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 404 {object} apiresponse.Response "No active subscription found"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/payment/subscription/active [get]
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
