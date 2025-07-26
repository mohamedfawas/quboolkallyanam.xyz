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

type subscriptionPlansRepository struct {
	db *postgres.Client
}

func NewSubscriptionPlansRepository(db *postgres.Client) repository.SubscriptionPlansRepository {
	return &subscriptionPlansRepository{db: db}
}

func (r *subscriptionPlansRepository) GetPlanByID(ctx context.Context, planID string) (*entity.SubscriptionPlan, error) {
	var plan entity.SubscriptionPlan
	if err := r.db.GormDB.WithContext(ctx).Where("id = ?", planID).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("GetPlanByID error in subscription plans repository: %v", err)
		return nil, err
	}
	return &plan, nil
}

func (r *subscriptionPlansRepository) CreatePlan(ctx context.Context, plan entity.SubscriptionPlan) error {
	if err := r.db.GormDB.WithContext(ctx).Create(&plan).Error; err != nil {
		log.Printf("CreatePlan error in subscription plans repository: %v", err)
		return err
	}
	return nil
}

func (r *subscriptionPlansRepository) UpdatePlan(ctx context.Context, planID string, plan entity.SubscriptionPlan) error {
	if err := r.db.GormDB.WithContext(ctx).Where("id = ?", planID).Updates(&plan).Error; err != nil {
		log.Printf("UpdatePlan error in subscription plans repository: %v", err)
		return err
	}
	return nil
}

func (r *subscriptionPlansRepository) GetActivePlans(ctx context.Context) ([]*entity.SubscriptionPlan, error) {
	var plans []*entity.SubscriptionPlan
	if err := r.db.GormDB.WithContext(ctx).Where("is_active = ?", true).Find(&plans).Error; err != nil {
		log.Printf("GetActivePlans error in subscription plans repository: %v", err)
		return nil, err
	}
	return plans, nil
}
