package payment

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.CreatePaymentOrderResponse, error) {
	paymentOrder, err := u.paymentClient.CreatePaymentOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	amountInINR := fmt.Sprintf("%s %.2f", paymentOrder.Currency, paymentOrder.Amount)

	paymentURL := fmt.Sprintf("%s/payment/checkout?order_id=%s", u.config.BaseURL, paymentOrder.OrderID)

	createPaymentOrderResponse := &dto.CreatePaymentOrderResponse{
		OrderID:         paymentOrder.OrderID,
		PaymentURL:      paymentURL,
		AmountInINR:     amountInINR,
		RazorpayOrderID: paymentOrder.RazorpayOrderID,
		RazorpayKeyID:   paymentOrder.RazorpayKeyID,
		PlanID:          paymentOrder.PlanID,
		ExpiresAt:       paymentOrder.ExpiresAt,
	}

	return createPaymentOrderResponse, nil
}
