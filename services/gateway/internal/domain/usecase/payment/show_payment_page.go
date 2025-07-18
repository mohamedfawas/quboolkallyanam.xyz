package payment

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
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

	gatewayURL := u.config.BaseURL + constants.HTTPHandlerVersionV1
	paymentPageResponse := &dto.ShowPaymentPageResponse{
		RazorpayOrderID:    response.RazorpayOrderID,
		RazorpayKeyID:      response.RazorpayKeyID,
		PlanID:             response.PlanID,
		Amount:             response.Amount,
		DisplayAmount:      response.DisplayAmount,
		PlanDurationInDays: response.PlanDurationInDays,
		GatewayURL:         gatewayURL,
	}
	return paymentPageResponse, nil
}
