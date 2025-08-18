package event

import (
	"context"

	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
)

type EventPublisher interface {
	PublishUserProfileUpdated(ctx context.Context, event userevents.UserProfileUpdatedEvent) error
	PublishUserInterestSent(ctx context.Context, event userevents.UserInterestSentEvent) error
	PublishMutualMatchCreated(ctx context.Context, event userevents.MutualMatchCreatedEvent) error
}
