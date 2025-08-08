package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/templates"
)

func (n *notificationUsecase) HandleUserAccountDeletion(ctx context.Context, userEmail string) error {
	return n.emailAdapter.SendEmail(ctx, model.EmailRequest{
		To:      userEmail,
		Subject: "Account Deletion Confirmation",
		Body:    templates.BuildAccountDeletionBody(userEmail),
	})
}