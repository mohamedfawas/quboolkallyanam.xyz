package payment

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
)

type paymentUsecase struct {
	paymentClient client.PaymentClient
	config        *config.Config
}

func NewPaymentUsecase(
	paymentClient client.PaymentClient,
	config *config.Config,
) *paymentUsecase {
	return &paymentUsecase{
		paymentClient: paymentClient,
		config:        config}
}
