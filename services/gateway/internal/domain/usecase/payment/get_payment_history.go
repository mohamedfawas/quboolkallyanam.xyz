package payment

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (p *paymentUsecase) GetPaymentHistory(ctx context.Context) ([]*dto.GetPaymentHistoryResponse, error) {
	paymentHistory, err := p.paymentClient.GetPaymentHistory(ctx)
	if err != nil {
		log.Printf("GetPaymentHistory error in payment usecase: %v", err)
		return nil, err
	}

	return paymentHistory, nil
}
