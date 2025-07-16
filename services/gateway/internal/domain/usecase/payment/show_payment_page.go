package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *paymentUsecase) ShowPaymentPage(ctx context.Context, razorpayOrderID string) (*dto.ShowPaymentPageResponse, error) {
	req := dto.ShowPaymentPageRequest{
		RazorpayOrderID: razorpayOrderID,
	}

	response, err := u.paymentClient.ShowPaymentPage(ctx, req)
	if err != nil {
		return nil, err
	}

	// Add gateway-specific data
	response.RazorpayOrderID = razorpayOrderID
	response.GatewayURL = u.config.BaseURL

	return response, nil
}
