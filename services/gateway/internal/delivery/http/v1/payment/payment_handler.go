package payment

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	paymentUsecase usecase.PaymentUsecase
	logger         *zap.Logger
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase, logger: logger}
}