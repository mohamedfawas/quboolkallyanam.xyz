package subscription

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error) {
	subscription, err := s.subscriptionsRepository.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return nil, apperrors.ErrActiveSubscriptionNotFound
	}

	return subscription, nil
}
