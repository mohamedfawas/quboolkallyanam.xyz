package payments

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (p *paymentUsecase) GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]*entity.GetPaymentHistoryResponse, error) {
	payments, err := p.paymentRepository.GetPaymentHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	var paymentHistory []*entity.GetPaymentHistoryResponse
	for _, payment := range payments {
		paymentHistory = append(paymentHistory, &entity.GetPaymentHistoryResponse{
			ID:              payment.ID,
			PlanID:          payment.PlanID,
			RazorpayOrderID: payment.RazorpayOrderID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			PaymentMethod:   payment.PaymentMethod,
			CreatedAt:       payment.CreatedAt,
		})
	}

	return paymentHistory, nil
}
