package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type AuthEventHandler struct {
	messagingClient    messageBroker.Client
	userProfileUsecase usecase.UserProfileUsecase
}

func NewAuthEventHandler(messagingClient messageBroker.Client, userProfileUsecase usecase.UserProfileUsecase) *AuthEventHandler {
	return &AuthEventHandler{
		messagingClient:    messagingClient,
		userProfileUsecase: userProfileUsecase,
	}
}

func (h *AuthEventHandler) StartListening(ctx context.Context) error {
	if h.messagingClient == nil {
		log.Println("messaging client is nil, skipping event listening for auth events")
		return nil
	}

	handler := func(data []byte) error {
		var authEvent authevents.UserLoginSuccessEvent
		if err := json.Unmarshal(data, &authEvent); err != nil {
			log.Printf("failed to unmarshal user login success event: %v", err)
			return err
		}

		return h.handleUserLoginSuccess(ctx, authEvent)
	}

	log.Println("starting to listen for user login success events")
	return h.messagingClient.Subscribe(constants.EventUserLoginSuccess, handler)
}

func (h *AuthEventHandler) handleUserLoginSuccess(ctx context.Context, event authevents.UserLoginSuccessEvent) error {
	log.Printf("[USER] Processing user login success event for user: %s", event.UserID)

	if err := h.userProfileUsecase.UpdateUserLastLogin(ctx, event.UserID, event.Email, event.Phone); err != nil {
		log.Printf("[USER] failed to update user last login: %v", err)
		return err
	}

	log.Printf("[USER] user last login updated successfully for user: %s", event.UserID)
	return nil
}
