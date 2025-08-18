package payment

import (
	"context"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) GetActiveSubscriptionPlans(ctx context.Context) ([]*dto.SubscriptionPlan, error) {
	plans, err := u.paymentClient.GetActiveSubscriptionPlans(ctx)
	if err != nil {
		return nil, err
	}

	return plans, nil
}