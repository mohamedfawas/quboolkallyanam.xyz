package usecase

import (
	"context"
	"strconv"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

func (n *notificationUsecase) HandleOTPVerification(ctx context.Context,
	userEmail, otp string, expiryMinutes int) error {

	expiryMinutesStr := strconv.Itoa(expiryMinutes)

	emailReq := model.EmailRequest{
		To:      userEmail,
		Subject: "OTP Verification",
		Body:    templates.BuildOTPVerificationBody(userEmail, otp, expiryMinutesStr),
	}

	return n.emailAdapter.SendEmail(ctx, emailReq)
}
