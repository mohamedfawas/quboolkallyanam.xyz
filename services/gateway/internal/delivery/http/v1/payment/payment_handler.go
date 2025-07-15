package payment

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type PaymentHandler struct {
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}
