package event

import (
	"context"
	"encoding/json"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	userEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/usecase"
	"go.uber.org/zap"
)

type EventHandler struct {
	messagingClient     messageBroker.Client
	notificationUsecase usecase.NotificationUsecase
	logger              *zap.Logger
}

func NewEventHandler(
	messagingClient messageBroker.Client,
	notificationUsecase usecase.NotificationUsecase,
	logger *zap.Logger,
) *EventHandler {
	return &EventHandler{
		messagingClient:     messagingClient,
		notificationUsecase: notificationUsecase,
		logger:              logger,
	}
}

func (h *EventHandler) StartListening(ctx context.Context) error {

	subscriptions := []struct {
		topic   string
		handler messageBroker.MessageHandler
	}{
		{
			topic:   constants.EventUserOTPRequested,
			handler: h.createUserOTPRequestedHandler(ctx),
		},
		{
			topic:   constants.EventUserAccountDeleted,
			handler: h.createUserAccountDeletionHandler(ctx),
		},
		{
			topic:   constants.EventAdminBlockedUser,
			handler: h.createAdminBlockedUserHandler(ctx),
		},
		{
			topic:   constants.EventUserInterestSent,
			handler: h.createUserInterestSentHandler(ctx),
		},
		{
			topic:   constants.EventMutualMatchCreated,
			handler: h.createMutualMatchCreatedHandler(ctx),
		},
	}

	for _, sub := range subscriptions {
		if err := h.messagingClient.Subscribe(sub.topic, sub.handler); err != nil {
			h.logger.Error("Error subscribing to topic", zap.String("topic", sub.topic), zap.Error(err))
			return err
		}
	}

	return nil
}

func (h *EventHandler) createUserOTPRequestedHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody authEvents.UserOTPRequestedEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			h.logger.Error("Error unmarshalling event", zap.Error(err))
			return err
		}

		return h.notificationUsecase.HandleOTPVerification(ctx, eventBody.Email, eventBody.OTP, eventBody.ExpiryMinutes)
	}
}

func (h *EventHandler) createUserAccountDeletionHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody authEvents.UserAccountDeletionEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			h.logger.Error("Error unmarshalling event", zap.Error(err))
			return err
		}

		return h.notificationUsecase.HandleUserAccountDeletion(ctx,  eventBody.Email)
	}
}

func (h *EventHandler) createAdminBlockedUserHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody authEvents.AdminBlockedUserEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			h.logger.Error("Error unmarshalling event", zap.Error(err))
			return err
		}

		return h.notificationUsecase.HandleAdminBlockedUser(ctx, eventBody.Email, eventBody.ShouldBlock)
	}
}

func (h *EventHandler) createUserInterestSentHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody userEvents.UserInterestSentEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			h.logger.Error("Error unmarshalling event", zap.Error(err))
			return err
		}

		return h.notificationUsecase.HandleUserInterestSent(ctx, eventBody.ReceiverEmail, eventBody.SenderProfileID, eventBody.SenderName)
	}
}

func (h *EventHandler) createMutualMatchCreatedHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody userEvents.MutualMatchCreatedEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			h.logger.Error("Error unmarshalling event", zap.Error(err))
			return err
		}

		return h.notificationUsecase.HandleMutualMatchCreated(ctx, 
			eventBody.User1Email, eventBody.User1ProfileID, eventBody.User1FullName, 
			eventBody.User2Email, eventBody.User2ProfileID, eventBody.User2FullName)
	}
}