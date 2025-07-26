package payment

import (
	"context"
	"log"

	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) GetSubscriptionPlan(ctx context.Context, planID string) (*dto.SubscriptionPlan, error) {
	if planID == "" {
		return nil, appError.ErrMissingRequiredFields
	}

	plan, err := u.paymentClient.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		log.Printf("GetSubscriptionPlan error in payment usecase: %v", err)
		return nil, err
	}

	if plan == nil {
		return nil, appError.ErrSubscriptionPlanNotFound
	}

	return plan, nil
}
