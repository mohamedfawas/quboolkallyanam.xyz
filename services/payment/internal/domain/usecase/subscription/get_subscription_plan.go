package subscription

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetSubscriptionPlan(ctx context.Context,
	planID string) (*entity.SubscriptionPlan, error) {
	if planID == "" {
		return nil, apperrors.ErrMissingRequiredFields
	}
	plan, err := s.subscriptionPlansRepository.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, apperrors.ErrSubscriptionPlanNotFound
	}

	return plan, nil
}
