package event

import (
	"context"

	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
)

type EventPublisher interface {
	PublishUserProfileUpdated(ctx context.Context, event userevents.UserProfileUpdatedEvent) error
}
