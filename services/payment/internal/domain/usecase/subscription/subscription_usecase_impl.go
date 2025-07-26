package subscription

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase"
)

type subscriptionUsecase struct {
	subscriptionPlansRepository repository.SubscriptionPlansRepository
	subscriptionsRepository     repository.SubscriptionsRepository
}

func NewSubscriptionUsecase(subscriptionPlansRepository repository.SubscriptionPlansRepository,
	subscriptionsRepository repository.SubscriptionsRepository) usecase.SubscriptionUsecase {
	return &subscriptionUsecase{
		subscriptionPlansRepository: subscriptionPlansRepository,
		subscriptionsRepository:     subscriptionsRepository,
	}
}
