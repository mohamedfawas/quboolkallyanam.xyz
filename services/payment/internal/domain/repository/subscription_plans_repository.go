package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type SubscriptionPlansRepository interface {
	GetPlanByID(ctx context.Context, planID string) (*entity.SubscriptionPlan, error)
}
