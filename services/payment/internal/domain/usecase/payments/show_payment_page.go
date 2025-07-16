package payments

import (
	"context"
	"errors"
	"fmt"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *paymentUsecase) ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*entity.ShowPaymentPageResponse, error) {
	// Get payment details
	payment, err := u.paymentRepository.GetPaymentDetailsByRazorpayOrderID(ctx, razorpayOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	// Get plan details
	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, payment.PlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrSubscriptionPlanNotFound
		}
		return nil, err
	}

	// Convert amount to paise (Razorpay expects amount in smallest currency unit)
	amountInPaise := int64(payment.Amount * 100)
	displayAmount := fmt.Sprintf("%.0f", payment.Amount)

	return &entity.ShowPaymentPageResponse{
		RazorpayOrderID:    payment.RazorpayOrderID,
		RazorpayKeyID:      u.razorpayService.KeyID(),
		PlanID:             plan.ID,
		Amount:             amountInPaise,
		DisplayAmount:      displayAmount,
		PlanDurationInDays: int32(plan.DurationDays),
	}, nil
}
