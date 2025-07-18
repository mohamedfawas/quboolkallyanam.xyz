package payments

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *paymentUsecase) ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*entity.ShowPaymentPageResponse, error) {
	payment, err := u.paymentRepository.GetPaymentDetailsByRazorpayOrderID(ctx, razorpayOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrPaymentNotFound
		}
		log.Printf("Error getting payment details by razorpay order id: %v", err)
		return nil, err
	}

	if payment.Status == constants.PaymentStatusCompleted {
		return nil, appErrors.ErrPaymentAlreadyCompleted
	}

	now := time.Now()
	if payment.ExpiresAt.Before(now) {
		log.Printf("current time: %v", now)
		log.Printf("Payment expired: %v", payment.ExpiresAt)
		return nil, appErrors.ErrPaymentExpired
	}

	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, payment.PlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrSubscriptionPlanNotFound
		}
		log.Printf("Error getting plan by id: %v", err)
		return nil, err
	}

	if !plan.IsActive {
		return nil, appErrors.ErrSubscriptionPlanNotActive
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
