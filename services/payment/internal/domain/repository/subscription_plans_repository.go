package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

type SubscriptionPlansRepository interface {
	GetPlanByID(ctx context.Context, planID string) (*entity.SubscriptionPlan, error)
	GetPlanByIDTx(ctx context.Context, tx *gorm.DB, planID string) (*entity.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, plan entity.SubscriptionPlan) error
	UpdatePlan(ctx context.Context, planID string, plan entity.SubscriptionPlan) error
	GetActivePlans(ctx context.Context) ([]*entity.SubscriptionPlan, error)
}