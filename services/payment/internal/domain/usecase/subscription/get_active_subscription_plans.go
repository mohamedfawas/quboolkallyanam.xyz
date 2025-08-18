package subscription

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetActiveSubscriptionPlans(ctx context.Context) ([]*entity.SubscriptionPlan, error) {
	plans, err := s.subscriptionPlansRepository.GetActivePlans(ctx)
	if err != nil {
		return nil, err
	}

	if len(plans) == 0 {
		return nil, apperrors.ErrSubscriptionPlanNotFound
	}

	return plans, nil
}
