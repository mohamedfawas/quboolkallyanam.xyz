package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

type SubscriptionsRepository interface {
	CreateSubscriptionTx(ctx context.Context, tx *gorm.DB, subscription *entity.Subscription) error
	GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error)
	GetActiveSubscriptionByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) (*entity.Subscription, error)
	UpdateSubscriptionTx(ctx context.Context, tx *gorm.DB, subscription *entity.Subscription) error
}