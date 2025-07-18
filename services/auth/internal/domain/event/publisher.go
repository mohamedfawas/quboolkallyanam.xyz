package event

import (
	"context"

	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
)

type EventPublisher interface {
	PublishUserOTPRequested(ctx context.Context, event authevents.UserOTPRequestedEvent) error
}
