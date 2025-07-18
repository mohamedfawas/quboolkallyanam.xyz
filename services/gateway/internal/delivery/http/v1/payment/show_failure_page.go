package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// @Summary Show failure page
// @Description Show failure page
// @Tags Payment
// @Accept json
// @Produce html
// @Param order_id query string true "Order ID"
// @Param error query string true "Error"
// @Success 200 {string} string "Success"
// @Router /api/v1/payment/failed [get]
func (h *PaymentHandler) ShowFailurePage(c *gin.Context) {
	failureData := dto.PaymentFailureData{
		OrderID: c.Query("order_id"),
		Error:   c.Query("error"),
	}

	c.HTML(http.StatusOK, "payment_failed.html", failureData)
}
