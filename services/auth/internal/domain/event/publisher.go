package event

import (
	"context"

	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
)

type EventPublisher interface {
	PublishUserOTPRequested(ctx context.Context, event authevents.UserOTPRequestedEvent) error
	PublishUserLoginSuccess(ctx context.Context, event authevents.UserLoginSuccessEvent) error
	PublishUserAccountDeletion(ctx context.Context, event authevents.UserAccountDeletionEvent) error
}
