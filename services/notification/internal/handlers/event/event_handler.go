package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/usecase"
)

type EventHandler struct {
	messagingClient     messageBroker.Client
	notificationUsecase usecase.NotificationUsecase
}

func NewEventHandler(
	messagingClient messageBroker.Client,
	notificationUsecase usecase.NotificationUsecase,
) *EventHandler {
	return &EventHandler{
		messagingClient:     messagingClient,
		notificationUsecase: notificationUsecase,
	}
}

func (h *EventHandler) StartListening(ctx context.Context) error {
	if h.messagingClient == nil {
		log.Println("Messaging client is nil, skipping event subscription")
		return nil
	}

	subscriptions := []struct {
		topic   string
		handler messageBroker.MessageHandler
	}{
		{
			topic:   constants.EventUserOTPRequested,
			handler: h.createUserOTPRequestedHandler(ctx),
		},
	}

	for _, sub := range subscriptions {
		if err := h.messagingClient.Subscribe(sub.topic, sub.handler); err != nil {
			log.Printf("Error subscribing to topic %s: %v", sub.topic, err)
			return err
		}
	}

	return nil
}

func (h *EventHandler) createUserOTPRequestedHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var eventBody authevents.UserOTPRequestedEvent

		if err := json.Unmarshal(body, &eventBody); err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			return err
		}

		return h.notificationUsecase.HandleOTPVerification(ctx, eventBody.Email, eventBody.OTP, eventBody.ExpiryMinutes)
	}
}
