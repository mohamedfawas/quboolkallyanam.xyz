package payment

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error) {
	verifyResp, err := u.paymentClient.VerifyPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}

	return verifyResp, nil
}
