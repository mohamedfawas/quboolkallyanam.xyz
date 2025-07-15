package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type SubscriptionsRepository interface {
	CreateSubscription(ctx context.Context, subscription *entity.Subscription) error
	GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *entity.Subscription) error
}
