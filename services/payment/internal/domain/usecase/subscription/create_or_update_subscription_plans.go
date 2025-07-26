package subscription

import (
	"context"
	"log"
	"time"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (s *subscriptionUsecase) CreateOrUpdateSubscriptionPlans(ctx context.Context,
	req entity.UpdateSubscriptionPlanRequest) error {

	if !validation.IsValidSubscriptionPlanID(req.ID) {
		return appError.ErrInvalidSubscriptionPlanID
	}

	existingPlan, err := s.subscriptionPlansRepository.GetPlanByID(ctx, req.ID)
	if err != nil {
		log.Printf("error getting subscription plan by id: %v", err)
		return err
	}

	if existingPlan == nil {
		if req.DurationDays == nil || req.Amount == nil || req.Currency == nil {
			return appError.ErrMissingRequiredFields
		}

		if !validation.IsValidSubscriptionPlanDurationDays(*req.DurationDays) {
			return appError.ErrInvalidSubscriptionPlanDurationDays
		}

		if !validation.IsValidSubscriptionPlanAmount(*req.Amount) {
			return appError.ErrInvalidSubscriptionPlanAmount
		}

		if !validation.IsValidSubscriptionPlanCurrency(*req.Currency) {
			return appError.ErrInvalidCurrency
		}

		description := ""
		if req.Description != nil {
			description = *req.Description
		}

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		now := time.Now().UTC()
		return s.subscriptionPlansRepository.CreatePlan(ctx, entity.SubscriptionPlan{
			ID:           req.ID,
			DurationDays: *req.DurationDays,
			Amount:       *req.Amount,
			Currency:     *req.Currency,
			Description:  description,
			IsActive:     isActive,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	if req.DurationDays != nil {
		if !validation.IsValidSubscriptionPlanDurationDays(*req.DurationDays) {
			return appError.ErrInvalidSubscriptionPlanDurationDays
		}
		existingPlan.DurationDays = *req.DurationDays
	}

	if req.Amount != nil {
		if !validation.IsValidSubscriptionPlanAmount(*req.Amount) {
			return appError.ErrInvalidSubscriptionPlanAmount
		}
		existingPlan.Amount = *req.Amount
	}

	if req.Currency != nil {
		if !validation.IsValidSubscriptionPlanCurrency(*req.Currency) {
			return appError.ErrInvalidCurrency
		}
		existingPlan.Currency = *req.Currency
	}

	if req.Description != nil {
		existingPlan.Description = *req.Description
	}

	if req.IsActive != nil {
		existingPlan.IsActive = *req.IsActive
	}

	existingPlan.UpdatedAt = time.Now().UTC()

	return s.subscriptionPlansRepository.UpdatePlan(ctx, req.ID, *existingPlan)
}
