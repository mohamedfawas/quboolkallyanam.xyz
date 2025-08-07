package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

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
