package subscription

import (
	"context"
	"log"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetActiveSubscriptionPlans(ctx context.Context) ([]*entity.SubscriptionPlan, error) {
	plans, err := s.subscriptionPlansRepository.GetActivePlans(ctx)
	if err != nil {
		log.Printf("GetActiveSubscriptionPlans error in subscription usecase: %v", err)
		return nil, err
	}

	if len(plans) == 0 {
		return nil, appError.ErrSubscriptionPlanNotFound
	}

	return plans, nil
}
