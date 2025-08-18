package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

func (n *notificationUsecase) HandleMutualMatchCreated(ctx context.Context,
	user1Email string, user1ProfileID int64, user1FullName string,
	user2Email string, user2ProfileID int64, user2FullName string) error {

	body := templates.BuildMutualMatchCommonBody(user1FullName, user1ProfileID, user2FullName, user2ProfileID)

	emailReq1 := model.EmailRequest{
		To:      user1Email,
		Subject: "New mutual match on Qubool Kallyanam",
		Body:    body,
	}

	emailReq2 := model.EmailRequest{
		To:      user2Email,
		Subject: "New mutual match on Qubool Kallyanam",
		Body:    body,
	}

	if err := n.emailAdapter.SendEmail(ctx, emailReq1); err != nil {
		return err
	}

	return n.emailAdapter.SendEmail(ctx, emailReq2)
}