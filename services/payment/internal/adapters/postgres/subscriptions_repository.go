package postgres

import (
	"context"
	"errors"
	"log"

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

func (r *subscriptionsRepository) CreateSubscription(ctx context.Context, subscription *entity.Subscription) error {
	db := GetDBFromContext(ctx, r.db.GormDB)
	return db.WithContext(ctx).Create(subscription).Error
}

func (r *subscriptionsRepository) GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error) {
	var subscription entity.Subscription
	db := GetDBFromContext(ctx, r.db.GormDB)

	err := db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND end_date > NOW()", userID, "active").
		First(&subscription).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("GetActiveSubscriptionByUserID error in subscriptions repository: %v", err)
		return nil, err
	}

	return &subscription, nil
}

func (r *subscriptionsRepository) UpdateSubscription(ctx context.Context, subscription *entity.Subscription) error {
	db := GetDBFromContext(ctx, r.db.GormDB)
	return db.WithContext(ctx).Save(subscription).Error
}
