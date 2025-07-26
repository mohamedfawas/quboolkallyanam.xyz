package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) CreateOrUpdateSubscriptionPlan(ctx context.Context, req dto.UpdateSubscriptionPlanRequest) (*dto.CreateOrUpdateSubscriptionPlanResponse, error) {
	response, err := u.paymentClient.CreateOrUpdateSubscriptionPlan(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
