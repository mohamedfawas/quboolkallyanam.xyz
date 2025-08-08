package usecase

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/port"
)

type notificationUsecase struct {
	emailAdapter port.EmailAdapter
}

func NewNotificationUsecase(emailAdapter port.EmailAdapter) NotificationUsecase {
	return &notificationUsecase{
		emailAdapter: emailAdapter,
	}
}
