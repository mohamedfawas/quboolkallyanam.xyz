package payment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *PaymentHandler) ShowPaymentPage(c *gin.Context) {
	razorpayOrderID := c.Query("razorpay_order_id")
	if razorpayOrderID == "" {
		// Use status found, otherwise it won't redirect correctly
		c.Redirect(http.StatusFound, "/payment/failed?error=missing_order_id")
		return
	}

	paymentPageResponseData, err := h.paymentUsecase.ShowPaymentPage(c.Request.Context(), razorpayOrderID)
	if err != nil {
		log.Printf("Failed to get payment page data: %v", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=failed_to_load", razorpayOrderID))
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentPageResponseData)
}
