package client

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type PaymentClient interface {
	CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.PaymentOrderResponse, error)
	ShowPaymentPage(ctx context.Context, req dto.ShowPaymentPageRequest) (*dto.ShowPaymentPageResponse, error)
	VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error)
	CreateOrUpdateSubscriptionPlan(ctx context.Context, req dto.UpdateSubscriptionPlanRequest) (*dto.CreateOrUpdateSubscriptionPlanResponse, error)
	GetSubscriptionPlan(ctx context.Context, planID string) (*dto.SubscriptionPlan, error)
	GetActiveSubscriptionPlans(ctx context.Context) ([]*dto.SubscriptionPlan, error)
}
