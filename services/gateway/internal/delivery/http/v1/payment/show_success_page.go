package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Payment success page
// @Description Render subscription success details
// @Tags Payment
// @Produce html
// @Param subscription_id query string true "Subscription ID"
// @Param start_date query string true "Start date (YYYY-MM-DD HH:MM:SS)"
// @Param end_date query string true "End date (YYYY-MM-DD HH:MM:SS)"
// @Param status query string true "Subscription status"
// @Success 200 {string} html "Success page"
// @Router /payment/success [get]
func (h *PaymentHandler) ShowSuccessPage(c *gin.Context) {
	subscriptionID := c.Query("subscription_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	successData := dto.PaymentSuccessData{
		SubscriptionID:        subscriptionID,
		SubscriptionStartDate: startDate,
		SubscriptionEndDate:   endDate,
		SubscriptionStatus:    status,
	}

	c.HTML(http.StatusOK, "payment_success.html", successData)
}
