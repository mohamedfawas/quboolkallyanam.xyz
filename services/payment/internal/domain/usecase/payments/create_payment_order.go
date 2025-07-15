package payments

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (u *paymentUsecase) CreatePaymentOrder(ctx context.Context, userID string, planID string) (*entity.PaymentOrderResponse, error) {
	plan, err := u.subscriptionPlanRepository.GetPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrSubscriptionPlanNotFound
		}
		return nil, err
	}

	if !plan.IsActive {
		return nil, appErrors.ErrSubscriptionPlanNotActive
	}

	razorpayOrderID, err := u.razorpayService.CreateOrder(plan.Amount, constants.PaymentCurrencyINR)
	if err != nil {
		return nil, err
	}

	orderID := "order_" + uuid.New().String()
	now := time.Now().UTC()
	expiresAt := now.Add(time.Minute * 15)

	payment := &entity.Payment{
		OrderID:           orderID,
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

	expiryTimeInIST := timeutil.ToIST(expiresAt)
	paymentOrderResponse := &entity.PaymentOrderResponse{
		OrderID:         orderID,
		RazorpayOrderID: razorpayOrderID,
		RazorpayKeyID:   u.razorpayService.KeyID(),
		Amount:          plan.Amount,
		Currency:        constants.PaymentCurrencyINR,
		PlanID:          planID,
		ExpiresAt:       expiryTimeInIST,
	}

	return paymentOrderResponse, nil
}
