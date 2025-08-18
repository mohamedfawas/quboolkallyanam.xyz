package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (u *paymentUsecase) ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*entity.ShowPaymentPageResponse, error) {
	payment, err := u.paymentRepository.GetPaymentDetailsByRazorpayOrderID(ctx, razorpayOrderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, apperrors.ErrPaymentNotFound
	}

	if payment.Status == constants.PaymentStatusCompleted {
		return nil, apperrors.ErrPaymentAlreadyCompleted
	}

	now := time.Now().UTC()
	if payment.ExpiresAt.Before(now) {
		return nil, apperrors.ErrPaymentExpired
	}

	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, payment.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, apperrors.ErrSubscriptionPlanNotFound
	}

	if !plan.IsActive {
		return nil, apperrors.ErrSubscriptionPlanNotActive
	}

	amountInPaise := int64(payment.Amount * 100)
	displayAmount := fmt.Sprintf("%s %.2f", constants.PaymentCurrencyINR, payment.Amount)

	return &entity.ShowPaymentPageResponse{
		RazorpayOrderID:    payment.RazorpayOrderID,
		RazorpayKeyID:      u.razorpayService.KeyID(),
		PlanID:             payment.PlanID,
		Amount:             amountInPaise,
		DisplayAmount:      displayAmount,
		PlanDurationInDays: int32(plan.DurationDays),
	}, nil
}
