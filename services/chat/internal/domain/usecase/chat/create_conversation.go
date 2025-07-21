package chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (c *chatUsecase) CreateConversation(ctx context.Context,
	userUUIDStr string, partnerProfileID int64) (*entity.Conversation, error) {

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		log.Printf("failed to parse user UUID: %v", err)
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	userProjection, err := c.userProjectionRepository.GetUserProjectionByID(ctx, partnerProfileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("user projection not found: %v", err)
			return nil, fmt.Errorf("user projection not found: %w", err)
		}
		log.Printf("failed to get user projection: %v", err)
		return nil, fmt.Errorf("failed to get user projection: %w", err)
	}

	partnerUUID, err := uuid.Parse(userProjection.UserUUID.String())
	if err != nil {
		log.Printf("failed to parse partner UUID: %v", err)
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
		log.Printf("failed to get conversation by participants: %v", err)
		return nil, fmt.Errorf("failed to get conversation by participants: %w", err)
	}

	if existingConversation != nil {
		log.Printf("conversation already exists: %v", existingConversation)
		return existingConversation, nil
	}

	now := time.Now().UTC()
	newConversation := &entity.Conversation{
		ID:             entity.NewConversationID(),
		ParticipantIDs: participantIDs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := c.conversationRepository.CreateConversation(ctx, newConversation); err != nil {
		log.Printf("failed to create conversation: %v", err)
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return newConversation, nil
}
