package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type PaymentUsecase interface {
	CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.CreatePaymentOrderResponse, error)
	ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*dto.ShowPaymentPageResponse, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error)
}
