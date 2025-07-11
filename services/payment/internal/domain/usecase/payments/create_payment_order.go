package payments

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (u *paymentUsecase) CreatePaymentOrder(ctx context.Context, userID string, planID string) (*entity.Payment, error) {
	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, errors.ErrSubscriptionPlanNotFound
	}

	if !plan.IsActive {
		return nil, errors.ErrSubscriptionPlanNotActive
	}

	razorpayOrderID, err := u.razorpayService.CreateOrder(plan.Amount, constants.PaymentCurrencyINR)
	if err != nil {
		return nil, err
	}

	now := timeutil.NowIST()
	payment := &entity.Payment{
		UserID:            userID,
		RazorpayOrderID:   razorpayOrderID,
		RazorpayPaymentID: "",
		RazorpaySignature: "",
		Amount:            plan.Amount,
		Currency:          constants.PaymentCurrencyINR,
		Status:            constants.PaymentStatusPending,
		PaymentMethod:     constants.PaymentMethodRazorpay,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := u.paymentRepository.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
