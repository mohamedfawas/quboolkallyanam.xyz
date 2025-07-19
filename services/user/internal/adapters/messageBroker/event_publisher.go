package messagebroker

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/event"
)

type eventPublisher struct {
	messagingClient messageBroker.Client
}

func NewEventPublisher(messagingClient messageBroker.Client) event.EventPublisher {
	return &eventPublisher{
		messagingClient: messagingClient,
	}
}

func (p *eventPublisher) PublishUserProfileCreated(
	ctx context.Context,
	event userevents.UserProfileCreatedEvent) error {

	if p.messagingClient == nil {
		log.Println("messaging client is nil, skipping event publishing UserProfileCreated event for user: %s", event.Email)
		return nil
	}

	if err := p.messagingClient.Publish(constants.EventUserProfileCreated, event); err != nil {
		log.Printf("failed to publish user profile created event for user: %s, error: %v", event.Email, err)
		return err
	}

	log.Printf("user profile created event published successfully for user: %s", event.Email)
	return nil
}
