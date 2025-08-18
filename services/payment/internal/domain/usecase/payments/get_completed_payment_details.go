package payments

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
)

func (u *paymentUsecase) GetCompletedPaymentDetails(ctx context.Context, page, limit int) ([]*entity.Payment, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return u.paymentRepository.GetCompletedPayments(ctx, limit, offset)
}