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

func (r *conversationRepository) GetConversationByID(ctx context.Context, conversationID string) (*entity.Conversation, error) {
	coll := r.db.Collection(constants.MongoDBCollectionConversations)

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, fmt.Errorf("mongodb: invalid conversation ID: %w", err)
	}

	filter := bson.M{"_id": objID}

	var conv entity.Conversation
	err = coll.FindOne(ctx, filter).Decode(&conv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("mongodb: find conversation: %w", err)
	}
	return &conv, nil
}

func (r *conversationRepository) GetUserConversations(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]*entity.Conversation, int64, error) {
	coll := r.db.Collection(constants.MongoDBCollectionConversations)

	// only conversations where userID is one of the participants
	filter := bson.M{"participant_ids": userID}

	totalCount, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("mongodb: count conversations: %w", err)
	}

	// build options: pagination + multi‑field sort
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{
			{"last_message_at", -1}, // most recent activity first
			{"created_at", -1},      // tie‑breaker: newest conversations first
		})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("mongodb: find conversations: %w", err)
	}
	defer cursor.Close(ctx)

	var conversations []*entity.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, 0, fmt.Errorf("mongodb: decode conversations: %w", err)
	}

	return conversations, totalCount, nil
}
