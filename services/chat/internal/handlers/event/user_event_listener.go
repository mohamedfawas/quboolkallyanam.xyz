package event

import (
	"context"
	"encoding/json"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"
	"go.uber.org/zap"
)

type UserEventListener struct {
	messagingClient       messageBroker.Client
	userProjectionUsecase usecase.UserProjectionUsecase
	logger                *zap.Logger
}

func NewUserEventListener(
	messagingClient messageBroker.Client,
	userProjectionUsecase usecase.UserProjectionUsecase,
	logger *zap.Logger) *UserEventListener {

	return &UserEventListener{
		messagingClient:       messagingClient,
		userProjectionUsecase: userProjectionUsecase,
		logger:                logger,
	}
}

func (h *UserEventListener) StartListening(
	ctx context.Context) error {
	handler := func(data []byte) error {
		var userEvent userevents.UserProfileUpdatedEvent
		if err := json.Unmarshal(data, &userEvent); err != nil {
			h.logger.Error("failed to unmarshal user profile updated event",
				zap.String(constants.UserIDS, userEvent.UserID.String()),
				zap.Error(err))
		}

		return h.handleUserProfileUpdated(ctx, userEvent)
	}

	h.logger.Info("starting to listen for user profile updated events")
	return h.messagingClient.Subscribe(constants.EventUserProfileUpdated, handler)
}

func (h *UserEventListener) handleUserProfileUpdated(ctx context.Context, event userevents.UserProfileUpdatedEvent) error {
	userProjection := &entity.UserProjection{
		UserUUID:      event.UserID,
		UserProfileID: event.UserProfileID,
		Email:         event.Email,
		FullName:      event.FullName,
	}

	if err := h.userProjectionUsecase.CreateOrUpdateUserProjection(ctx, userProjection); err != nil {
		h.logger.Error("failed to create user projection",
			zap.String(constants.UserIDS, event.UserID.String()),
			zap.Error(err))
		return err
	}

	return nil
}
