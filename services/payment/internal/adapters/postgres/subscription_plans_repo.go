package postgres

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
)

type subscriptionPlansRepository struct {
	db *postgres.Client
}

func NewSubscriptionPlansRepository(db *postgres.Client) repository.SubscriptionPlansRepository {
	return &subscriptionPlansRepository{db: db}
}

func (r *subscriptionPlansRepository) GetPlanByID(ctx context.Context, planID string) (*entity.SubscriptionPlan, error) {
	var plan entity.SubscriptionPlan
	if err := r.db.GormDB.WithContext(ctx).Where("id = ?", planID).First(&plan).Error; err != nil {
		log.Printf("GetPlanByID error in subscription plans repository: %v", err)
		return nil, err
	}
	return &plan, nil
}
