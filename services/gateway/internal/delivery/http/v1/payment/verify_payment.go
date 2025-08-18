package payment

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/metrics"
	"go.uber.org/zap"
)

// @Summary Verify Razorpay payment
// @Description Verify Razorpay signature and redirect to success/failure page
// @Tags Payment
// @Produce html
// @Param razorpay_order_id query string true "Razorpay order ID"
// @Param razorpay_payment_id query string true "Razorpay payment ID"
// @Param razorpay_signature query string true "Razorpay signature"
// @Success 302 {string} string "Redirect to success page"
// @Failure 302 {string} string "Redirect to failure page"
// @Router /payment/verify [get]
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract request context", zap.Error(err))
		c.Redirect(http.StatusFound, "/payment/failed?error=context_error")
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	req := dto.VerifyPaymentRequest{
		RazorpayOrderID:   c.Query("razorpay_order_id"),
		RazorpayPaymentID: c.Query("razorpay_payment_id"),
		RazorpaySignature: c.Query("razorpay_signature"),
	}
	if req.RazorpayOrderID == "" || req.RazorpayPaymentID == "" || req.RazorpaySignature == "" {
		c.Redirect(http.StatusFound, "/payment/failed?error=missing_parameters")
		return
	}

	verifyResp, err := h.paymentUsecase.VerifyPayment(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Payment verification failed", zap.Error(err), zap.String("razorpay_order_id", req.RazorpayOrderID))
		}
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

	metrics.PaymentsCompletedTotal.Inc()
	log.Info("Payment verified successfully", zap.String("subscription_id", verifyResp.SubscriptionID))
	c.Redirect(http.StatusFound, successURL)
}
