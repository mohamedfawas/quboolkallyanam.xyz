package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)


// @Summary Payment failure page
// @Description Render failure details
// @Tags Payment
// @Produce html
// @Param order_id query string false "Order ID"
// @Param error query string false "Error code"
// @Success 200 {string} html "Failure page"
// @Router /payment/failed [get]
func (h *PaymentHandler) ShowFailurePage(c *gin.Context) {
	failureData := dto.PaymentFailureData{
		OrderID: c.Query("order_id"),
		Error:   c.Query("error"),
	}

	c.HTML(http.StatusOK, "payment_failed.html", failureData)
}
