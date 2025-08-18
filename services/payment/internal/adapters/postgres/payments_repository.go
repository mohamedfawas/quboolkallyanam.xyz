package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
	"gorm.io/gorm"
)

type paymentsRepository struct {
	db *postgres.Client
}

func NewPaymentsRepository(db *postgres.Client) repository.PaymentsRepository {
	return &paymentsRepository{db: db}
}

func (r *paymentsRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	return r.db.GormDB.WithContext(ctx).Create(payment).Error
}

func (r *paymentsRepository) GetPaymentDetailsByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.GormDB.WithContext(ctx).
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentsRepository) GetPaymentDetailsByRazorpayOrderIDTx(ctx context.Context, tx *gorm.DB, razorpayOrderID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := tx.WithContext(ctx).
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentsRepository) UpdatePaymentTx(ctx context.Context, tx *gorm.DB, payment *entity.Payment) error {
	return tx.WithContext(ctx).Save(payment).Error
}

func (r *paymentsRepository) GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.db.GormDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}


func (r *paymentsRepository) GetCompletedPayments(ctx context.Context, limit, offset int) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.db.GormDB.WithContext(ctx).
		Where("status = ?", constants.PaymentStatusCompleted).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}