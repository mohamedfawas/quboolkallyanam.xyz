package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
)

type paymentsRepository struct {
	db *postgres.Client
}

func NewPaymentsRepository(db *postgres.Client) repository.PaymentsRepository {
	return &paymentsRepository{db: db}
}

func (r *paymentsRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	db := GetDBFromContext(ctx, r.db.GormDB)
	return db.WithContext(ctx).Create(payment).Error
}

func (r *paymentsRepository) GetPaymentDetailsByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (*entity.Payment, error) {
	var payment entity.Payment
	db := GetDBFromContext(ctx, r.db.GormDB)

	err := db.WithContext(ctx).
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentsRepository) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	db := GetDBFromContext(ctx, r.db.GormDB)
	return db.WithContext(ctx).Save(payment).Error
}

func (r *paymentsRepository) GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	db := GetDBFromContext(ctx, r.db.GormDB)

	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}
