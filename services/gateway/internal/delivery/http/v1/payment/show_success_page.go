package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *PaymentHandler) ShowSuccessPage(c *gin.Context) {
	successData := dto.PaymentSuccessData{
		OrderID:    c.Query("order_id"),
		PaymentID:  c.Query("payment_id"),
		Amount:     c.Query("amount"),
		ValidUntil: c.Query("valid_until"),
	}

	c.HTML(http.StatusOK, "payment_success.html", successData)
}
