package client

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type PaymentClient interface {
	CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.PaymentOrderResponse, error)
}
