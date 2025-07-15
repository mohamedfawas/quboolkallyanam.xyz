package v1

import (
	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// /////////////////////////// Create Payment Order //////////////////////////////
func MapCreatePaymentOrderRequest(req dto.PaymentOrderRequest) *paymentpbv1.CreatePaymentOrderRequest {
	return &paymentpbv1.CreatePaymentOrderRequest{
		PlanId: req.PlanID,
	}
}

func MapCreatePaymentOrderResponse(resp *paymentpbv1.CreatePaymentOrderResponse) *dto.PaymentOrderResponse {
	return &dto.PaymentOrderResponse{
		OrderID:         resp.OrderId,
		RazorpayOrderID: resp.RazorpayOrderId,
		RazorpayKeyID:   resp.RazorpayKeyId,
		Amount:          resp.Amount,
		Currency:        resp.Currency,
		PlanID:          resp.PlanId,
		ExpiresAt:       resp.ExpiresAt.AsTime(),
	}
}
