package usecase

import (
	"context"
	"strconv"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/port"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

type notificationUsecase struct {
	emailAdapter port.EmailAdapter
}

func NewNotificationUsecase(emailAdapter port.EmailAdapter) NotificationUsecase {
	return &notificationUsecase{
		emailAdapter: emailAdapter,
	}
}

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
