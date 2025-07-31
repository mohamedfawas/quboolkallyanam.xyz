package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type PaymentsRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentDetailsByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (*entity.Payment, error)
	UpdatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.Payment, error)
}
