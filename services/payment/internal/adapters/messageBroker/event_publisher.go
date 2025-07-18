package messagebroker

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/event"
)

type eventPublisher struct {
	messagingClient messageBroker.Client
}

func NewEventPublisher(messagingClient messageBroker.Client) event.EventPublisher {
	return &eventPublisher{
		messagingClient: messagingClient,
	}
}

func (p *eventPublisher) PublishPaymentVerified(ctx context.Context,
	event paymentEvents.PaymentVerified) error {

	if p.messagingClient == nil {
		log.Println("messaging client is nil, skipping event publishing")
		return nil
	}

	if err := p.messagingClient.Publish(constants.EventUserPaymentVerified, event); err != nil {
		log.Printf("failed to publish payment verified event: %v", err)
		return err
	}

	log.Printf("payment verified event published successfully for user: %s", event.UserID)
	return nil
}
