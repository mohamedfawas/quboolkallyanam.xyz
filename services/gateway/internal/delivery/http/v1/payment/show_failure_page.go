package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *PaymentHandler) ShowFailurePage(c *gin.Context) {
	failureData := dto.PaymentFailureData{
		OrderID: c.Query("order_id"),
		Error:   c.Query("error"),
	}

	c.HTML(http.StatusOK, "payment_failed.html", failureData)
}
