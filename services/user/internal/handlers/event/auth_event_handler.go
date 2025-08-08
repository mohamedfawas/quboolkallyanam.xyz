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
	handler := func(data []byte) error {
		var authEvent authevents.UserLoginSuccessEvent
		if err := json.Unmarshal(data, &authEvent); err != nil {
			h.logger.Error("failed to unmarshal user login success event",
				zap.String(constants.UserIDS, authEvent.UserID.String()),
				zap.Error(err))
			return err
		}

		return h.handleUserLoginSuccess(ctx, authEvent)
	}

	h.logger.Info("starting to listen for user login success events")
	return h.messagingClient.Subscribe(constants.EventUserLoginSuccess, handler)
}

func (h *AuthEventHandler) handleUserLoginSuccess(ctx context.Context, event authevents.UserLoginSuccessEvent) error {

	if err := h.userProfileUsecase.UpdateUserLastLogin(ctx, event.UserID, event.Email, event.Phone); err != nil {
		h.logger.Error("failed to update user last login",
			zap.String(constants.UserIDS, event.UserID.String()),
			zap.Error(err))
		return err
	}

	h.logger.Info("user last login updated successfully",
		zap.String(constants.UserIDS, event.UserID.String()))
	return nil
}
