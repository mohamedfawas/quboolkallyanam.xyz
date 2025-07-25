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

// @Summary Create payment order
// @Description Create payment order
// @Tags Payment
// @Accept json
// @Produce json
// @Param payment_order_request body dto.PaymentOrderRequest true "Payment order request"
// @Success 200 {object} dto.PaymentOrderResponse "Payment order response"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/payment/order [post]
func (h *PaymentHandler) CreatePaymentOrder(c *gin.Context) {
	var req dto.PaymentOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		log.Printf("User ID not found in context")
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	response, err := h.paymentUsecase.CreatePaymentOrder(ctx, req)
	if err != nil {
		log.Printf("Failed to create payment order: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Payment order created successfully", response)
}
