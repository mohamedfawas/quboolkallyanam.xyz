package chat

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *chatUsecase) CreateConversation(ctx context.Context,
	userUUIDStr string, partnerProfileID int64) (*entity.CreateConversationResponse, error) {

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	userProjection, err := c.userProjectionRepository.GetUserProjectionByUUID(ctx, userUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user projection: %w", err)
	}
	if userProjection == nil {
		return nil, appError.ErrUserNotFound
	}

	partnerUserProjection, err := c.userProjectionRepository.GetUserProjectionByID(ctx, partnerProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user projection: %w", err)
	}

	if partnerUserProjection == nil {
		return nil, appError.ErrPartnerUserNotFound
	}

	// For creating response we need names
	partnerName := partnerUserProjection.FullName
	userName := userProjection.FullName
	participantNames := []string{userName, partnerName}

	partnerUUID, err := uuid.Parse(partnerUserProjection.UserUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse partner UUID: %w", err)
	}

	participants := []string{userUUID.String(), partnerUUID.String()}
	sort.Strings(participants)

	// Convert []string to []entity.UserID
	participantIDs := make([]entity.UserID, len(participants))
	for i, p := range participants {
		participantIDs[i] = entity.UserID(p)
	}

	existingConversation, err := c.conversationRepository.GetConversationByParticipants(ctx, participants)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to get conversation by participants: %w", err)
	}

	if existingConversation != nil {
		return &entity.CreateConversationResponse{
			ConversationID: existingConversation.ID,
			Participants:   participantNames,
			CreatedAt:      existingConversation.CreatedAt,
		}, nil
	}

	now := time.Now().UTC()
	newConversation := &entity.Conversation{
		ID:             entity.NewConversationID(),
		ParticipantIDs: participantIDs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := c.conversationRepository.CreateConversation(ctx, newConversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return &entity.CreateConversationResponse{
		ConversationID: newConversation.ID,
		Participants:   participantNames,
		CreatedAt:      newConversation.CreatedAt,
	}, nil
}
