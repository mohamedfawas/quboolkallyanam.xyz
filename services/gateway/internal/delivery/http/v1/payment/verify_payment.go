package payment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)


func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		RazorpayOrderID   string `form:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `form:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `form:"razorpay_signature" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Invalid verification request: %v", err)
		c.Redirect(http.StatusFound, "/payment/failed?error=invalid_request")
		return
	}

	verifyReq := dto.VerifyPaymentRequest{
		RazorpayOrderID:   req.RazorpayOrderID,
		RazorpayPaymentID: req.RazorpayPaymentID,
		RazorpaySignature: req.RazorpaySignature,
	}

	verifyResp, err := h.paymentUsecase.VerifyPayment(c.Request.Context(), verifyReq)
	if err != nil {
		log.Printf("Payment verification failed: %v", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=%s", 
			req.RazorpayOrderID, "verification_failed"))
		return
	}

	// Redirect to success page with payment details
	c.Redirect(http.StatusFound, fmt.Sprintf("/payment/success?order_id=%s&payment_id=%s&amount=%s&valid_until=%s", 
		req.RazorpayOrderID, req.RazorpayPaymentID, "999", verifyResp.SubscriptionEndDate))
}