package payment

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Verify payment
// @Description Verify payment
// @Tags Payment
// @Accept json
// @Produce html
// @Param razorpay_order_id query string true "Razorpay order ID"
// @Param razorpay_payment_id query string true "Razorpay payment ID"
// @Param razorpay_signature query string true "Razorpay signature"
// @Success 200 {string} string "Success"
// @Router /api/v1/payment/verify [get]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {

	req := dto.VerifyPaymentRequest{
		RazorpayOrderID:   c.Query("razorpay_order_id"),
		RazorpayPaymentID: c.Query("razorpay_payment_id"),
		RazorpaySignature: c.Query("razorpay_signature"),
	}

	if req.RazorpayOrderID == "" || req.RazorpayPaymentID == "" || req.RazorpaySignature == "" {
		log.Printf("Missing required verification parameters")
		c.Redirect(http.StatusFound, "/payment/failed?error=missing_parameters")
		return
	}

	verifyResp, err := h.paymentUsecase.VerifyPayment(c.Request.Context(), req)
	if err != nil {
		log.Printf("Payment verification failed: %v", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=%s",
			req.RazorpayOrderID, "verification_failed"))
		return
	}

	startDate := verifyResp.SubscriptionStartDate.Format("2006-01-02 15:04:05")
	endDate := verifyResp.SubscriptionEndDate.Format("2006-01-02 15:04:05")

	successURL := fmt.Sprintf("/payment/success?subscription_id=%s&start_date=%s&end_date=%s&status=%s",
		url.QueryEscape(verifyResp.SubscriptionID),
		url.QueryEscape(startDate),
		url.QueryEscape(endDate),
		url.QueryEscape(verifyResp.SubscriptionStatus))

	c.Redirect(http.StatusFound, successURL)
}
