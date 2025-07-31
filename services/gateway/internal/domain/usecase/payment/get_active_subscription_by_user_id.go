package payment

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (p *paymentUsecase) GetActiveSubscriptionByUserID(ctx context.Context) (*dto.ActiveSubscription, error) {
	activeSubscription, err := p.paymentClient.GetActiveSubscriptionByUserID(ctx)
	if err != nil {
		log.Printf("GetActiveSubscriptionByUserID error in payment usecase: %v", err)
		return nil, err
	}

	return activeSubscription, nil
}
