package messagebroker

import (
	"context"

	"go.uber.org/zap"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/event"
)

type eventPublisher struct {
	messagingClient messageBroker.Client
	logger          *zap.Logger
}

func NewEventPublisher(messagingClient messageBroker.Client, logger *zap.Logger) event.EventPublisher {
	return &eventPublisher{
		messagingClient: messagingClient,
		logger:          logger,
	}
}

func (p *eventPublisher) PublishPaymentVerified(ctx context.Context,
	event paymentEvents.PaymentVerified) error {

	if err := p.messagingClient.Publish(constants.EventUserPaymentVerified, event); err != nil {
		p.logger.Error("failed to publish payment verified event", zap.Error(err))
		return err
	}

	p.logger.Info("payment verified event published successfully for user", zap.String("user_id", event.UserID))
	return nil
}
