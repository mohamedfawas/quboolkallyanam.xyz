package usecase	

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

func (n *notificationUsecase) HandleAdminBlockedUser(ctx context.Context,
	userEmail string, shouldBlock bool) error {

	emailReq := model.EmailRequest{
		To:      userEmail,
		Subject: "Account reviewed by an admin",
		Body:    templates.BuildAdminBlockOrUnblockUserStatusBody(userEmail, shouldBlock),
	}

	return n.emailAdapter.SendEmail(ctx, emailReq)
}