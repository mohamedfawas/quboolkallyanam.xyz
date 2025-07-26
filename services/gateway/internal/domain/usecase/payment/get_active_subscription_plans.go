package payment

import (
	"context"
	"log"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) GetActiveSubscriptionPlans(ctx context.Context) ([]*dto.SubscriptionPlan, error) {
	plans, err := u.paymentClient.GetActiveSubscriptionPlans(ctx)
	if err != nil {
		log.Printf("GetActiveSubscriptionPlans error in payment usecase: %v", err)
		return nil, err
	}

	return plans, nil
}