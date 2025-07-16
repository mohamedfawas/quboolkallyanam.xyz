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
		RazorpayOrderID: resp.RazorpayOrderId,
		Amount:          resp.Amount,
		Currency:        resp.Currency,
		PlanID:          resp.PlanId,
		ExpiresAt:       resp.ExpiresAt.AsTime(),
	}
}

// /////////////////////////// Show Payment Page //////////////////////////////
func MapShowPaymentPageRequest(req dto.ShowPaymentPageRequest) *paymentpbv1.ShowPaymentPageRequest {
	return &paymentpbv1.ShowPaymentPageRequest{
		RazorpayOrderId: req.RazorpayOrderID,
	}
}

func MapShowPaymentPageResponse(resp *paymentpbv1.ShowPaymentPageResponse) *dto.ShowPaymentPageResponse {
	return &dto.ShowPaymentPageResponse{
		RazorpayOrderID:    "", // Will be set in usecase
		RazorpayKeyID:      "", // Will be set in usecase
		PlanID:             resp.PlanId,
		Amount:             0, // Will be set in usecase
		DisplayAmount:      resp.DisplayAmount,
		PlanDurationInDays: int(resp.PlanDurationInDays),
		GatewayURL:         "", // Will be set in usecase
	}
}

// /////////////////////////// Verify Payment //////////////////////////////
func MapVerifyPaymentRequest(req dto.VerifyPaymentRequest) *paymentpbv1.VerifyPaymentRequest {
	return &paymentpbv1.VerifyPaymentRequest{
		RazorpayOrderId:   req.RazorpayOrderID,
		RazorpayPaymentId: req.RazorpayPaymentID,
		RazorpaySignature: req.RazorpaySignature,
	}
}

func MapVerifyPaymentResponse(resp *paymentpbv1.VerifyPaymentResponse) *dto.VerifyPaymentResponse {
	return &dto.VerifyPaymentResponse{
		SubscriptionID:        resp.SubscriptionId,
		SubscriptionStartDate: resp.SubscriptionStartDate.AsTime().Format("2006-01-02 15:04:05"),
		SubscriptionEndDate:   resp.SubscriptionEndDate.AsTime().Format("2006-01-02 15:04:05"),
		SubscriptionStatus:    resp.SubscriptionStatus,
	}
}
