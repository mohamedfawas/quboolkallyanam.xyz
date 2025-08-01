package mongodb

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(client *mongodb.Client) repository.MessageRepository {
	return &messageRepository{
		db: client.Database,
	}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *entity.Message) error {
	coll := r.db.Collection(constants.MongoDBCollectionMessages)
	if _, err := coll.InsertOne(ctx, message); err != nil {
		return fmt.Errorf("mongodb: insert message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetMessagesByConversationID(ctx context.Context, conversationID string, limit, offset int32) ([]*entity.Message, int64, error) {
	coll := r.db.Collection(constants.MongoDBCollectionMessages)

	conversationObjID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, 0, fmt.Errorf("mongodb: invalid conversation ID: %w", err)
	}

	filter := bson.M{"conversation_id": conversationObjID}

	// Get total count
	totalCount, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("mongodb: count messages: %w", err)
	}

	// Setup pagination options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "sent_at", Value: -1}}) // Sort by newest first
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("mongodb: find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*entity.Message
	for cursor.Next(ctx) {
		var message entity.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, 0, fmt.Errorf("mongodb: decode message: %w", err)
		}
		messages = append(messages, &message)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("mongodb: cursor error: %w", err)
	}

	return messages, totalCount, nil
}
