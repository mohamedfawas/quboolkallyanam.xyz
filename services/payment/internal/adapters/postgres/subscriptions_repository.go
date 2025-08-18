package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
)

type subscriptionsRepository struct {
	db *postgres.Client
}

func NewSubscriptionsRepository(db *postgres.Client) repository.SubscriptionsRepository {
	return &subscriptionsRepository{db: db}
}

func (r *subscriptionsRepository) CreateSubscriptionTx(ctx context.Context, tx *gorm.DB, subscription *entity.Subscription) error {
	return tx.WithContext(ctx).Create(subscription).Error
}

func (r *subscriptionsRepository) GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.GormDB.WithContext(ctx).
		Where("user_id = ? AND status = ? AND end_date > NOW()", userID, "active").
		First(&subscription).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &subscription, nil
}

func (r *subscriptionsRepository) GetActiveSubscriptionByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := tx.WithContext(ctx).
		Where("user_id = ? AND status = ? AND end_date > NOW()", userID, "active").
		First(&subscription).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &subscription, nil
}

func (r *subscriptionsRepository) UpdateSubscriptionTx(ctx context.Context, tx *gorm.DB, subscription *entity.Subscription) error {
	return tx.WithContext(ctx).Save(subscription).Error
}