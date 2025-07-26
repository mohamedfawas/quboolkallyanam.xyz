package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type SubscriptionUsecase interface {
	CreateOrUpdateSubscriptionPlans(ctx context.Context,
		req entity.UpdateSubscriptionPlanRequest) error
	GetSubscriptionPlan(ctx context.Context,
		planID string) (*entity.SubscriptionPlan, error)
	GetActiveSubscriptionPlans(ctx context.Context) ([]*entity.SubscriptionPlan, error)
	GetActiveSubscriptionByUserID(ctx context.Context,
		userID string) (*entity.Subscription, error)
}
