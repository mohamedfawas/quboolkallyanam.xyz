package payment

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.CreatePaymentOrderResponse, error) {
	paymentOrder, err := u.paymentClient.CreatePaymentOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	paymentURL := fmt.Sprintf("%s/payment/checkout?razorpay_order_id=%s", u.config.BaseURL, paymentOrder.RazorpayOrderID)
	displayAmount := fmt.Sprintf("%s %.2f", constants.PaymentCurrencyINR, paymentOrder.Amount)

	return &dto.CreatePaymentOrderResponse{
		PaymentURL:      paymentURL,
		RazorpayOrderID: paymentOrder.RazorpayOrderID,
		Amount:          displayAmount,
		PlanID:          paymentOrder.PlanID,
		ExpiresAt:       paymentOrder.ExpiresAt,
	}, nil
}
