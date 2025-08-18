package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) GetSubscriptionPlan(ctx context.Context, planID string) (*dto.SubscriptionPlan, error) {
	if planID == "" {
		return nil, apperrors.ErrMissingRequiredFields
	}

	plan, err := u.paymentClient.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
