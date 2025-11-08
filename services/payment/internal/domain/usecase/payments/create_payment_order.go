package payments

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (u *paymentUsecase) CreatePaymentOrder(ctx context.Context, userID string, planID string) (*entity.PaymentOrderResponse, error) {
	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, apperrors.ErrSubscriptionPlanNotFound
	}

	if !plan.IsActive {
		return nil, apperrors.ErrSubscriptionPlanNotActive
	}

	razorpayOrderID, err := u.razorpayService.CreateOrder(plan.Amount, constants.PaymentCurrencyINR)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Minute * 30)

	payment := &entity.Payment{
		UserID:            userID,
		PlanID:            planID,
		RazorpayOrderID:   razorpayOrderID,
		RazorpayPaymentID: "",
		RazorpaySignature: "",
		Amount:            plan.Amount,
		Currency:          constants.PaymentCurrencyINR,
		Status:            constants.PaymentStatusPending,
		PaymentMethod:     constants.PaymentMethodRazorpay,
		ExpiresAt:         expiresAt,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := u.paymentRepository.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	paymentOrderResponse := &entity.PaymentOrderResponse{
		RazorpayOrderID: razorpayOrderID,
		Amount:          plan.Amount,
		Currency:        constants.PaymentCurrencyINR,
		PlanID:          planID,
		ExpiresAt:       expiresAt,
	}

	return paymentOrderResponse, nil
}
