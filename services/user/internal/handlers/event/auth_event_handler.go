package event

import (
	"context"
	"encoding/json"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
	"go.uber.org/zap"
)

type AuthEventHandler struct {
	messagingClient    messageBroker.Client
	userProfileUsecase usecase.UserProfileUsecase
	logger             *zap.Logger
}

func NewAuthEventHandler(
	messagingClient messageBroker.Client,
	userProfileUsecase usecase.UserProfileUsecase,
	logger *zap.Logger) *AuthEventHandler {

	return &AuthEventHandler{
		messagingClient:    messagingClient,
		userProfileUsecase: userProfileUsecase,
		logger:             logger,
	}
}

func (h *AuthEventHandler) StartListening(ctx context.Context) error {
	subscriptions := []struct {
		topic   string
		handler messageBroker.MessageHandler
	}{
		{
			topic:   constants.EventUserLoginSuccess,
			handler: h.createUserLoginSuccessHandler(ctx),
		},
		{
			topic:   constants.EventUserAccountDeleted,
			handler: h.createUserAccountDeletionHandler(ctx),
		},
		// Add other events here
	}

	for _, sub := range subscriptions {
		if err := h.messagingClient.Subscribe(sub.topic, sub.handler); err != nil {
			h.logger.Error("Error subscribing to topic",
				zap.String("topic", sub.topic),
				zap.Error(err))
			return err
		}
		h.logger.Info("Successfully subscribed to topic", zap.String("topic", sub.topic))
	}

	return nil
}

func (h *AuthEventHandler) createUserLoginSuccessHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var event authevents.UserLoginSuccessEvent

		if err := json.Unmarshal(body, &event); err != nil {
			h.logger.Error("Error unmarshalling user login success event", zap.Error(err))
			return err
		}

		if err := h.userProfileUsecase.UpdateUserLastLogin(ctx, event.UserID, event.Email, event.Phone); err != nil {
			h.logger.Error("failed to update user last login",
				zap.String(constants.ContextKeyUserID, event.UserID.String()),
				zap.Error(err))
			return err
		}

		h.logger.Info("user last login updated successfully",
			zap.String(constants.ContextKeyUserID, event.UserID.String()))
		return nil
	}
}

func (h *AuthEventHandler) createUserAccountDeletionHandler(ctx context.Context) messageBroker.MessageHandler {
	return func(body []byte) error {
		var event authevents.UserAccountDeletionEvent

		if err := json.Unmarshal(body, &event); err != nil {
			h.logger.Error("Error unmarshalling user account deletion event", zap.Error(err))
			return err
		}

		if err := h.userProfileUsecase.HandleUserDeletion(ctx, event.UserID); err != nil {
			h.logger.Error("failed to handle user deletion",
				zap.String(constants.ContextKeyUserID, event.UserID.String()),
				zap.Error(err))
			return err
		}

		h.logger.Info("user account deletion handled successfully",
			zap.String(constants.ContextKeyUserID, event.UserID.String()))
		return nil
	}
}
