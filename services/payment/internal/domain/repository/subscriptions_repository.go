package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type SubscriptionsRepository interface {
	GetSubscriptions(ctx context.Context, userID string) ([]*entity.Subscription, error)
	GetSubscriptionByID(ctx context.Context, subscriptionID string) (*entity.Subscription, error)
	CreateSubscription(ctx context.Context, subscription *entity.Subscription) error
	UpdateSubscription(ctx context.Context, subscriptionID string, subscription *entity.Subscription) error
}
