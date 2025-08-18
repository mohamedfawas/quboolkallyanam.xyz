package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) GetCompletedPaymentDetails(
	ctx context.Context,
	req dto.GetCompletedPaymentDetailsRequest,
) (*dto.GetCompletedPaymentDetailsResponse, error) {
	if req.Page < 1 {
		return nil, apperrors.ErrInvalidPaginationPage
	}
	if req.Limit < 1 || req.Limit > 50 {
		return nil, apperrors.ErrInvalidPaginationLimit
	}

	return u.paymentClient.GetCompletedPaymentDetails(ctx, req)
}