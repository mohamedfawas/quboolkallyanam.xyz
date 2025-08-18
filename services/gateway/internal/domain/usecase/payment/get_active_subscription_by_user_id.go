package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (p *paymentUsecase) GetActiveSubscriptionByUserID(ctx context.Context) (*dto.ActiveSubscription, error) {
	activeSubscription, err := p.paymentClient.GetActiveSubscriptionByUserID(ctx)
	if err != nil {
		return nil, err
	}

	return activeSubscription, nil
}
