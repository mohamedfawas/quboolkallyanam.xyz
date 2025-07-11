package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type PaymentsRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentByID(ctx context.Context, paymentID string) (*entity.Payment, error)
	UpdatePayment(ctx context.Context, paymentID string, payment *entity.Payment) error
}
