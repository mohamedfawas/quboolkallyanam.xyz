package event

import (
	"context"
	"encoding/json"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
	"go.uber.org/zap"
)

type PaymentEventHandler struct {
	messagingClient messageBroker.Client
	userUsecase     usecase.UserUsecase
	logger          *zap.Logger
}

func NewPaymentEventHandler(
	messagingClient messageBroker.Client,
	userUsecase usecase.UserUsecase,
	logger *zap.Logger,
) *PaymentEventHandler {
	return &PaymentEventHandler{
		messagingClient: messagingClient,
		userUsecase:     userUsecase,
		logger:          logger,
	}
}

func (h *PaymentEventHandler) StartListening(ctx context.Context) error {
	handler := func(data []byte) error {
		var paymentEvent paymentEvents.PaymentVerified
		if err := json.Unmarshal(data, &paymentEvent); err != nil {
			h.logger.Error("Failed to unmarshal payment verified event", zap.Error(err))
			return err
		}

		return h.handlePaymentVerified(ctx, paymentEvent)
	}

	h.logger.Info("Starting to listen for payment verified events")
	return h.messagingClient.Subscribe(constants.EventUserPaymentVerified, handler)
}

func (h *PaymentEventHandler) handlePaymentVerified(ctx context.Context, event paymentEvents.PaymentVerified) error {
	h.logger.Info("Processing payment verified event for user", zap.String("user_id", event.UserID))

	if err := h.userUsecase.UpdateUserPremium(ctx, event.UserID, event.SubscriptionEndDate); err != nil {
		h.logger.Error("Failed to update user premium status", zap.Error(err))
		return err
	}

	h.logger.Info("Successfully processed payment verified event for user", zap.String("user_id", event.UserID))
	return nil
}
