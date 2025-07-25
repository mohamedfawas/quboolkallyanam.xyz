package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type PaymentUsecase interface {
	CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.CreatePaymentOrderResponse, error)
	ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*dto.ShowPaymentPageResponse, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error)
	CreateOrUpdateSubscriptionPlan(ctx context.Context, req dto.UpdateSubscriptionPlanRequest) (*dto.CreateOrUpdateSubscriptionPlanResponse, error)
	GetSubscriptionPlan(ctx context.Context, planID string) (*dto.SubscriptionPlan, error)
	GetActiveSubscriptionPlans(ctx context.Context) ([]*dto.SubscriptionPlan, error)
}
