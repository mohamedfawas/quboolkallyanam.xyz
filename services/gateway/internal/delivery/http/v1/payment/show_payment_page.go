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
		c.Redirect(http.StatusBadRequest, "/payment/failed?error=missing_order_id")
		return
	}

	paymentData, err := h.paymentUsecase.ShowPaymentPage(c.Request.Context(), razorpayOrderID)
	if err != nil {
		log.Printf("Failed to get payment page data: %v", err)
		c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=failed_to_load", razorpayOrderID))
		return
	}

	c.HTML(http.StatusOK, "payment.html", paymentData)
}
