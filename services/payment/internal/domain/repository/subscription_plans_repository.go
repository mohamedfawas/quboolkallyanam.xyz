package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type SubscriptionPlansRepository interface {
	GetPlanByID(ctx context.Context, planID string) (*entity.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, plan entity.SubscriptionPlan) error
	UpdatePlan(ctx context.Context, planID string, plan entity.SubscriptionPlan) error
	GetActivePlans(ctx context.Context) ([]*entity.SubscriptionPlan, error)
}
