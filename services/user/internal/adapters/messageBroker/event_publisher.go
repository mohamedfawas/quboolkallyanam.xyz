package messagebroker

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/event"
	"go.uber.org/zap"
)

type eventPublisher struct {
	messagingClient messageBroker.Client
	logger          *zap.Logger
}

func NewEventPublisher(
	messagingClient messageBroker.Client,
	logger *zap.Logger) event.EventPublisher {

	return &eventPublisher{
		messagingClient: messagingClient,
		logger:          logger,
	}
}

func (p *eventPublisher) PublishUserProfileUpdated(
	ctx context.Context,
	event userevents.UserProfileUpdatedEvent) error {

	if err := p.messagingClient.Publish(constants.EventUserProfileUpdated, event); err != nil {
		p.logger.Error("failed to publish user profile updated event",
			zap.String(constants.UserIDS, event.UserID.String()),
			zap.Error(err))
		return err
	}

	return nil
}

func (p *eventPublisher) PublishUserInterestSent(
	ctx context.Context,
	event userevents.UserInterestSentEvent) error {

	if err := p.messagingClient.Publish(constants.EventUserInterestSent, event); err != nil {
		p.logger.Error("failed to publish user interest sent event",
			zap.String(constants.UserIDS, event.ReceiverEmail),
			zap.Error(err))
		return err
	}
	return nil
}

func (p *eventPublisher) PublishMutualMatchCreated(
	ctx context.Context,
	event userevents.MutualMatchCreatedEvent) error {

	if err := p.messagingClient.Publish(constants.EventMutualMatchCreated, event); err != nil {
		p.logger.Error("failed to publish mutual match created event",
			zap.String(constants.UserIDS, event.User1Email),
			zap.String(constants.UserIDS, event.User2Email),
			zap.Error(err))
		return err
	}
	return nil
}
