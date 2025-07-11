package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type PaymentUsecase interface {
	CreatePaymentOrder(ctx context.Context, userID string, planID string) (*entity.Payment, error)
	// VerifyPayment(ctx context.Context, orderID string, paymentID string, signature string) error
}
