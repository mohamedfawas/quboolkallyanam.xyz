package payment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Show payment checkout page
// @Description Render Razorpay checkout page for an order
// @Tags Payment
// @Produce html
// @Param razorpay_order_id query string true "Razorpay order ID"
// @Success 200 {string} html "Payment page"
// @Failure 302 {string} string "Redirected on error"
// @Router /payment/checkout [get]
func (h *PaymentHandler) ShowPaymentPage(c *gin.Context) {
	// Log BEFORE context extraction to catch all requests
	h.logger.Info("ShowPaymentPage called",
		zap.String("razorpay_order_id", c.Query("razorpay_order_id")),
		zap.String("path", c.Request.URL.Path))
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract request context", zap.Error(err))
		c.Redirect(http.StatusFound, "/payment/failed?error=context_error")
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	razorpayOrderID := c.Query("razorpay_order_id")
	if razorpayOrderID == "" {
		c.Redirect(http.StatusFound, "/payment/failed?error=missing_order_id")
		return
	}

	resp, err := h.paymentUsecase.ShowPaymentPage(reqCtx.Ctx, razorpayOrderID)
	if err != nil {
		// Log ALL errors
		log.Error("Failed to get payment page data",
			zap.Error(err),
			zap.String("razorpay_order_id", razorpayOrderID))

		c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=failed_to_load", razorpayOrderID))
		return

		/*
			if apperrors.ShouldLogError(err) {
				log.Error("Failed to get payment page data", zap.Error(err), zap.String("razorpay_order_id", razorpayOrderID))
			}
			c.Redirect(http.StatusFound, fmt.Sprintf("/payment/failed?order_id=%s&error=failed_to_load", razorpayOrderID))
			return
		*/
	}

	log.Info("Payment page data is retrieved", zap.String("razorpay_order_id", razorpayOrderID))
	c.HTML(http.StatusOK, "payment.html", resp)
}
