package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

type PaymentUsecase interface {
	CreatePaymentOrder(ctx context.Context, userID string, planID string) (*entity.PaymentOrderResponse, error)
	ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*entity.ShowPaymentPageResponse, error)
	VerifyPayment(ctx context.Context, req *entity.VerifyPaymentRequest) (*entity.VerifyPaymentResponse, error)
	GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.GetPaymentHistoryResponse, error)
}
