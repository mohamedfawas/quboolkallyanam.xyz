package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type PaymentEventHandler struct {
	messagingClient messageBroker.Client
	userUsecase     usecase.UserUsecase
}

func NewPaymentEventHandler(messagingClient messageBroker.Client, userUsecase usecase.UserUsecase) *PaymentEventHandler {
	return &PaymentEventHandler{
		messagingClient: messagingClient,
		userUsecase:     userUsecase,
	}
}

func (h *PaymentEventHandler) StartListening(ctx context.Context) error {
	if h.messagingClient == nil {
		log.Println("Messaging client is nil, skipping event subscription")
		return nil
	}

	handler := func(data []byte) error {
		var paymentEvent paymentEvents.PaymentVerified
		if err := json.Unmarshal(data, &paymentEvent); err != nil {
			log.Printf("Failed to unmarshal payment verified event: %v", err)
			return err
		}

		return h.handlePaymentVerified(ctx, paymentEvent)
	}

	log.Println("Starting to listen for payment verified events")
	return h.messagingClient.Subscribe(constants.EventUserPaymentVerified, handler)
}

func (h *PaymentEventHandler) handlePaymentVerified(ctx context.Context, event paymentEvents.PaymentVerified) error {
	log.Printf("[AUTH] Processing payment verified event for user: %s", event.UserID)

	if err := h.userUsecase.UpdateUserPremium(ctx, event.UserID, event.SubscriptionEndDate); err != nil {
		log.Printf("[AUTH] Failed to update user premium status: %v", err)
		return err
	}

	log.Printf("[AUTH] Successfully processed payment verified event for user: %s", event.UserID)
	return nil
}
