package event

import (
	"context"

	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
)

type EventPublisher interface {
	PublishUserProfileCreated(ctx context.Context, event userevents.UserProfileCreatedEvent) error
}
