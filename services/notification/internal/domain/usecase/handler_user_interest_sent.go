package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

func (n *notificationUsecase) HandleUserInterestSent(ctx context.Context,
	receiverEmail string, senderProfileID int64, senderName string) error {

	emailReq := model.EmailRequest{
		To:      receiverEmail,
		Subject: "New interest on Qubool Kallyanam",
		Body:    templates.BuildInterestSentBody(receiverEmail, senderProfileID, senderName),
	}

	return n.emailAdapter.SendEmail(ctx, emailReq)
}