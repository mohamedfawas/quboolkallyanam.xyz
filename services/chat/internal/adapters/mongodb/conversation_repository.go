package mongodb

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type conversationRepository struct {
	db *mongo.Database
}

func NewConversationRepository(client *mongodb.Client) repository.ConversationRepository {
	return &conversationRepository{
		db: client.Database,
	}
}

func (r *conversationRepository) CreateConversation(ctx context.Context, conv *entity.Conversation) error {
	coll := r.db.Collection(constants.MongoDBCollectionConversations)
	if _, err := coll.InsertOne(ctx, conv); err != nil {
		return fmt.Errorf("mongodb: insert conversation: %w", err)
	}
	return nil
}

func (r *conversationRepository) GetConversationByParticipants(ctx context.Context, participants []string) (*entity.Conversation, error) {
	coll := r.db.Collection(constants.MongoDBCollectionConversations)

	filter := bson.M{
		"participant_ids": bson.M{
			"$all":  participants,
			"$size": len(participants),
		},
	}

	var conv entity.Conversation
	err := coll.FindOne(ctx, filter).Decode(&conv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("mongodb: find conversation: %w", err)
	}
	return &conv, nil
}
