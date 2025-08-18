package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

type PaymentsRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentDetailsByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (*entity.Payment, error)
	GetPaymentDetailsByRazorpayOrderIDTx(ctx context.Context, tx *gorm.DB, razorpayOrderID string) (*entity.Payment, error)
	UpdatePaymentTx(ctx context.Context, tx *gorm.DB, payment *entity.Payment) error
	GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.Payment, error)
	GetCompletedPayments(ctx context.Context, limit, offset int) ([]*entity.Payment, error)
}