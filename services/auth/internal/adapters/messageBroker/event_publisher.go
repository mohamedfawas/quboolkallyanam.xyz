package messagebroker

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/event"
)

type eventPublisher struct {
	messagingClient messageBroker.Client
}

func NewEventPublisher(messagingClient messageBroker.Client) event.EventPublisher {
	return &eventPublisher{
		messagingClient: messagingClient,
	}
}

func (p *eventPublisher) PublishUserOTPRequested(ctx context.Context,
	event authevents.UserOTPRequestedEvent) error {

	if p.messagingClient == nil {
		log.Println("messaging client is nil, skipping event publishing UserOTPRequested event for user: %s", event.Email)
		return nil
	}

	if err := p.messagingClient.Publish(constants.EventUserOTPRequested, event); err != nil {
		log.Printf("failed to publish user otp requested event for user: %s, error: %v", event.Email, err)
		return err
	}

	log.Printf("user otp requested event published successfully for user: %s", event.Email)
	return nil
}
