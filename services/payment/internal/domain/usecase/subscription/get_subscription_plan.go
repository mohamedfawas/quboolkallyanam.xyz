package subscription

import (
	"context"

	"log"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetSubscriptionPlan(ctx context.Context,
	planID string) (*entity.SubscriptionPlan, error) {
	if planID == "" {
		return nil, appError.ErrMissingRequiredFields
	}
	plan, err := s.subscriptionPlansRepository.GetPlanByID(ctx, planID)
	if err != nil {
		log.Printf("GetSubscriptionPlan error in subscription usecase: %v", err)
		return nil, err
	}

	if plan == nil {
		return nil, appError.ErrSubscriptionPlanNotFound
	}

	return plan, nil
}
