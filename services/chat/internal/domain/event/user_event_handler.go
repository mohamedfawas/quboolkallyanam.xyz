package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"
)

type UserEventHandler struct {
	messagingClient messageBroker.Client
	userUsecase     usecase.UserProjectionUsecase
}

func NewUserEventHandler(messagingClient messageBroker.Client, userUsecase usecase.UserProjectionUsecase) *UserEventHandler {
	return &UserEventHandler{
		messagingClient: messagingClient,
		userUsecase:     userUsecase,
	}
}

func (h *UserEventHandler) StartListening(ctx context.Context) error {
	if h.messagingClient == nil {
		log.Println("Messaging client is nil, skipping event subscription")
		return nil
	}

	handler := func(data []byte) error {
		var userEvent userevents.UserProfileCreatedEvent
		if err := json.Unmarshal(data, &userEvent); err != nil {
			log.Printf("Failed to unmarshal user profile created event: %v", err)
			return err
		}

		return h.handleUserProfileCreated(ctx, userEvent)
	}

	log.Println("Starting to listen for user profile created events")
	return h.messagingClient.Subscribe(constants.EventUserProfileCreated, handler)
}

func (h *UserEventHandler) handleUserProfileCreated(ctx context.Context, event userevents.UserProfileCreatedEvent) error {
	log.Printf("[CHAT] Processing user profile created event for user: %s", event.UserID)

	if err := h.userUsecase.CreateUserProjection(ctx, &entity.UserProjection{
		UserUUID:      event.UserID,
		UserProfileID: event.UserProfileID,
		Email:         event.Email,
		Phone:         event.Phone,
		CreatedAt:     event.CreatedAt,
		UpdatedAt:     event.UpdatedAt,
	}); err != nil {
		log.Printf("[CHAT] Failed to create user projection: %v", err)
		return err
	}

	log.Printf("[CHAT] Successfully processed user profile created event for user: %s", event.UserID)
	return nil
}
