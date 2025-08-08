package messagebroker

import (
	"context"

	"go.uber.org/zap"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/event"
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

func (p *eventPublisher) PublishUserOTPRequested(ctx context.Context,
	event authevents.UserOTPRequestedEvent) error {
	if err := p.messagingClient.Publish(constants.EventUserOTPRequested, event); err != nil {
		p.logger.Error("failed to publish user otp requested event for user", zap.String("user_email", event.Email), zap.Error(err))
		return err
	}

	p.logger.Info("user otp requested event published successfully for user", zap.String("user_email", event.Email))
	return nil
}

func (p *eventPublisher) PublishUserLoginSuccess(ctx context.Context,
	event authevents.UserLoginSuccessEvent) error {

	if err := p.messagingClient.Publish(constants.EventUserLoginSuccess, event); err != nil {
		p.logger.Error("failed to publish user last login event for user", zap.String("user_email", event.Email), zap.Error(err))
		return err
	}

	p.logger.Info("user last login event published successfully for user", zap.String("user_email", event.Email))
	return nil
}

func (p *eventPublisher) PublishUserAccountDeletion(ctx context.Context,
	event authevents.UserAccountDeletionEvent) error {

	if err := p.messagingClient.Publish(constants.EventUserAccountDeleted, event); err != nil {
		p.logger.Error("failed to publish user account deletion event for user", zap.String("user_email", event.Email), zap.Error(err))
		return err
	}

	p.logger.Info("user account deletion event published successfully for user", zap.String("user_email", event.Email))
	return nil
}
