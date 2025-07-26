package subscription

import (
	"context"
	"log"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*entity.Subscription, error) {
	subscription, err := s.subscriptionsRepository.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetActiveSubscriptionByUserID error in subscription usecase: %v", err)
		return nil, err
	}

	if subscription == nil {
		return nil, appError.ErrActiveSubscriptionNotFound
	}

	return subscription, nil
}
